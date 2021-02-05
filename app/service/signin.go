package service

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/EventBus"
	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/crypto"

	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/oauth2"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/passwd"
	"github.com/suisrc/zgo/modules/store"

	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
)

// Signin 账户管理
type Signin struct {
	gpa.GPA                          // 数据库句柄
	Passwd         *passwd.Validator // 密码验证器
	Store          store.Storer      // 缓存控制器
	MSender        MobileSender      // 手机发送验证
	ESender        EmailSender       // 邮箱发送验证
	TSender        ThreeSender       // 三方发送验证
	OAuth2Selector oauth2.Selector   // OAuth2选择器
	Bus            EventBus.Bus      // 时间总线

}

// Signin 登陆控制
// params: c 访问上下文
// params: b 请求参数
// params: l 验证最后一次登录结果
// result: 登录者信息
func (a *Signin) Signin(c *gin.Context, b *schema.SigninBody, l func(*gin.Context, int64) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	if b.Password != "" {
		// 使用密码方式登录
		return a.SigninByPasswd(c, b, l)
	}
	if b.Captcha != "" {
		// 使用验证码方式登录
		if b.Code == "" {
			// 没有签名密钥， 验证码无效
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TYPE-CODE", Other: "校验码无效"})
		}
		// 执行验证
		return a.SigninByCaptcha(c, b, l)
	}
	// 没有合理的登录方式， 无法登录
	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TYPE-NONE", Other: "无效登陆方式"})
}

//============================================================================================

// SigninByPasswd 密码登陆
func (a *Signin) SigninByPasswd(c *gin.Context, b *schema.SigninBody, last func(*gin.Context, int64) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	// 查询账户信息
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, b.Username, schema.AccountTypeName, b.Platform, b.OrgCode, true); err != nil || account.ID <= 0 {
		// 无法查询到账户， 是否可以使用 2(手机)， 3(邮箱) 查询， 待定
		account.ID = 0
		if p, _ := regexp.MatchString(`^(1[3-8]\d{9}$`, b.Username); p {
			// 使用手机方式登录(只匹配中国手机号)
			err = account.QueryByParentAccount(a.Sqlx, b.Username, schema.AccountTypeMobile, b.Platform, b.OrgCode)
		} else if p, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+$`, b.Username); p {
			// 使用邮箱方式登录
			err = account.QueryByParentAccount(a.Sqlx, b.Username, schema.AccountTypeEmail, b.Platform, b.OrgCode)
		}
		if err != nil || account.ID <= 0 {
			// 登录失败， 最终无法完成登录的账户查询
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
		}
	}
	// 验证密码
	if b, err := a.VerifyPassword(b.Password, &account); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 密码验证发生异常
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	} else if !b {
		// 密码不匹配
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}
	if last != nil {
		if _, err := last(c, account.ID); err != nil {
			return nil, err // 快速验证上次登录结果
		}
	}
	// 获取用户信息
	return a.GetSignUserInfo(c, b, &account)
}

// SigninByCaptcha 验证码登陆
func (a *Signin) SigninByCaptcha(c *gin.Context, b *schema.SigninBody, last func(*gin.Context, int64) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	// 查询账户信息
	accountID, captchaGetter, err := DecryptCaptchaByAccount(c, b.Code)
	if err != nil {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-CODE", Other: "校验码不正确"})
	}
	// 查询账户信息
	account := schema.SigninGpaAccount{}
	if err := account.QueryByID(a.Sqlx, accountID); err != nil || account.ID <= 0 {
		logger.Errorf(c, logger.ErrorWW(err)) // 账户异常
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
	} else if account.Account != b.Username || account.PlatformKID.String != b.Platform || !account.VerifySecret.Valid {
		// 无法处理， 登录时候， 账户发生了变更
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
	}
	// 验证验证码
	if captcha, expire, err := captchaGetter(account.VerifySecret.String); err != nil {
		return nil, err // 解密验证码发生异常
	} else if expire <= 0 {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-EXPIRED", Other: "验证码已过期"})
	} else if captcha != b.Captcha {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-CHECK", Other: "验证码不正确"})
	}
	if last != nil {
		if _, err := last(c, account.ID); err != nil {
			return nil, err // 快速验证上次登录结果
		}
	}
	// 获取用户信息
	return a.GetSignUserInfo(c, b, &account)
}

//============================================================================================

// GetSignUserInfo with role
func (a *Signin) GetSignUserInfo(c *gin.Context, b *schema.SigninBody, sa *schema.SigninGpaAccount) (*schema.SigninUser, error) {
	if sa.Status != schema.StatusEnable { // 账户被禁用
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ACCOUNT-DISABLE", Other: "账户被禁用,请联系管理员"})
	}
	// 登陆用户
	suser := schema.SigninUser{}
	suser.Scope = b.Scope
	//suser.AccountID = strconv.Itoa(sa.ID) // SigninUser -> 1
	//suser.UserIdxID = strconv.Itoa(sa.UserID)
	suser.TokenID = jwt.NewTokenID(strconv.Itoa(int(sa.ID + 1103)))
	suser.Account, _ = EncryptAccountWithUser(c, sa.ID, sa.UserID.Int64, suser.TokenID) // 账户信息
	suser.TokenPID, _ = helper.GetCtxValueToString(c, helper.ResTknKey)                 // 子母令牌
	suser.Account2 = sa.CustomID.String

	if err := a.SetSignUserWithUser(c, sa, &suser); err != nil { // 用户信息
		return nil, err
	}
	if err := a.SetSignUserWithToken(c, b, sa, &suser); err != nil { // 访问令牌签名
		return nil, err
	}
	if suser.OrgAdmin != schema.SuperUser {
		// 如果是超级管理员， 需要跳过所有认证
		if err := a.SetSignUserWithRole(c, sa, &suser); err != nil { // 角色信息
			return nil, err
		}
	}

	return &suser, nil
}

// SetSignUserWithRole with role info
// 查询用户的所有的角色
// 如果一个人具有管理员权限， 其所有的角色都会被舍弃， 只保留管理员角色
func (a *Signin) SetSignUserWithRole(c *gin.Context, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	// 如果账户上带有角色， 优先使用账户角色登录系统Account
	if roles, err := new(schema.SigninGpaAccountRole).QueryAllByUserAndOrg(a.Sqlx, sa.ID, suser.OrgCode); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err // 数据库发生异常
		}
	} else if len(*roles) > 0 {
		// 处理得到的用户角色列表
		rs := []string{}
		for _, r := range *roles {
			if r.OrgAdm {
				// 一旦用户具有管理员角色， 系统会无视其他所有角色的使用
				suser.OrgAdmin = schema.SuperUser
				return nil
			} else if r.SvcCode.Valid {
				rs = append(rs, r.SvcCode.String+":"+r.Name) // 应用角色
			} else {
				rs = append(rs, r.Name) // 租户角色
			}
			rs = append(rs, helper.IfString(r.SvcCode.Valid, r.SvcCode.String+":"+r.Name, r.Name))
		}
		// 设定用户角色
		suser.SetUserRoles(rs)
		return nil // 使用了账户上的角色登录系统
	}
	if !sa.UserID.Valid {
		return nil // 账户上没有用户信息， 结束处理， 该账户无角色信息
	}
	// 账户上没有角色， 取用户在对应租户下的所有角色
	if roles, err := new(schema.SigninGpaUserRole).QueryAllByUserAndOrg(a.Sqlx, sa.UserID.Int64, suser.OrgCode); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err // 数据库发生异常
		}
	} else if len(*roles) > 0 {
		// 处理得到的用户角色列表
		rs := []string{}
		for _, r := range *roles {
			if r.OrgAdm {
				// 一旦用户具有管理员角色， 系统会无视其他所有角色的使用
				suser.OrgAdmin = schema.SuperUser
				return nil
			} else if r.SvcCode.Valid {
				rs = append(rs, r.SvcCode.String+":"+r.Name) // 应用角色
			} else {
				rs = append(rs, r.Name) // 租户角色
			}
		}
		// 设定用户角色
		suser.SetUserRoles(rs)
		return nil // 使用了用户上的角色登录系统
	}
	return nil // 账户和用户上都没有角色
}

// SetSignUserWithToken with token
// 登录客户端加密方式
func (a *Signin) SetSignUserWithToken(c *gin.Context, b *schema.SigninBody, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	cgw := schema.ClientGpaWebToken{}
	if b.WebToken != "" {
		// 使用指定的令牌
		if err := cgw.QueryByKID(a.Sqlx, b.WebToken); err != nil {
			// 令牌没有
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-WEB-TOKEN-NONE", Other: "JWT令牌密钥不存在"})
		} else if cgw.Status != schema.StatusPrivate {
			// 令牌状态
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-WEB-TOKEN-INVALID", Other: "JWT令牌密钥失效"})
		} else if cgw.OrgCode.Valid && (cgw.Type.String != "org" || cgw.OrgCode.String != suser.OrgCode) {
			// 专用令牌， 无权使用
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-WEB-TOKEN-UNAUTH", Other: "JWT令牌密钥无权使用"})
		} else if cgw.Target.Valid && (cgw.Type.String != "usr" || cgw.Target.Int64 != sa.UserID.Int64) {
			// 专用令牌， 无权使用
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-WEB-TOKEN-UNAUTH", Other: "JWT令牌密钥无权使用"})
		}
	}
	// 验证捆绑性质
	//  else if err := cgw.QueryByUsr(a.Sqlx, sa.UserID.Int64); err == nil {
	// 	// 判断当前用户是否需要专用JWT
	// 	// do nothings
	// } else if err := cgw.QueryByOrg(a.Sqlx, b.Org); err == nil {
	// 	// 判断当前租户是否需要专用JWT
	// 	// do nothings
	// }
	if cgw.KID != "" {
		helper.SetCtxValue(c, helper.ResJwtKey, cgw.KID) // 配置客户端, 该内容会影响JWT签名方式
		if cgw.JwtIssuer.Valid {
			suser.Issuer = cgw.JwtIssuer.String
		}
		if cgw.JwtAudience.Valid {
			suser.Audience = cgw.JwtAudience.String
		}
	}
	if suser.Issuer == "" {
		suser.Issuer = c.Request.Host
	}
	if suser.Audience == "" {
		suser.Audience = c.Request.Host
	}
	return nil
}

// SetSignUserWithUser with user info
// 查询用户信息
func (a *Signin) SetSignUserWithUser(c *gin.Context, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	if !sa.UserID.Valid {
		// 账户上没有用户信息， 待验证账户， 允许登录
		suser.OrgCode = sa.OrgCode.String
		if len(sa.Account) > 16 {
			suser.UserName = sa.Account[:16] + "..."
		} else {
			suser.UserName = sa.Account
		}
		return nil
	}
	if sa.OrgCode.Valid {
		// 账户上绑定了租户， 使用用户的租户账户
		user := schema.SigninGpaOrgUser{}
		if err := user.QueryByUserAndOrg(a.Sqlx, sa.UserID.Int64, sa.OrgCode.String); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是租户用户不存在
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户机构成员"})
		} else if user.Status != schema.StatusEnable {
			// 租户账户被禁用
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
		}
		suser.UserID = user.UnionKID
		suser.UserName = user.Name
		suser.OrgCode = user.OrgCode
		suser.OrgUsrID = user.CustomID.String
		if user.Type == schema.ORG {
			suser.OrgAdmin = schema.SuperUser // 租户根账户即是super user
		}
	} else {
		// 账户上未绑定租户， 使用用户的平台账户
		user := schema.SigninGpaUser{}
		if err := user.QueryByID(a.Sqlx, sa.UserID.Int64, ""); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户信息异常"})
		} else if user.Status != schema.StatusEnable {
			// 平台账户被禁用
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
		}
		suser.UserID = user.KID
		suser.UserName = user.Name
		suser.OrgCode = ""
		suser.OrgUsrID = ""
		if user.ID == 1 {
			suser.OrgCode = schema.PlatformCode // 修正平台编码
			suser.OrgAdmin = schema.SuperUser   // 平台超级管理员账户
		}
	}
	return nil
}

//============================================================================================

// VerifyPassword 验证密码
func (a *Signin) VerifyPassword(pwd string, acc *schema.SigninGpaAccount) (bool, error) {
	ok, _ := a.Passwd.Verify(&PasswdCheck{
		Account:  acc,
		Password: pwd,
	})
	return ok, nil
}

// PasswdCheck 密码验证实体
type PasswdCheck struct {
	Account  *schema.SigninGpaAccount
	Password string
}

var _ passwd.IEntity = &PasswdCheck{}

// Target 输入的密码
func (a *PasswdCheck) Target() string {
	return a.Password
}

// Source 保存的加密密码
func (a *PasswdCheck) Source() string {
	return a.Account.Password.String
}

// Salt 密码盐值
func (a *PasswdCheck) Salt() string {
	return a.Account.PasswordSalt.String
}

// Type 加密类型
func (a *PasswdCheck) Type() string {
	return a.Account.PasswordType.String
}

//============================================================================================

// Captcha 发送验证码
func (a *Signin) Captcha(c *gin.Context, b *schema.SigninOfCaptcha, k string) (string, error) {

	info, err := a.parseCaptchaType(c, b)
	if err != nil {
		return "", err
	}
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, info.Account, info.Type, info.Platform, b.OrgCode, true); err != nil {
		if !sqlxc.IsNotFound(err) {
			return "", err
		}
		return "", helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
	}
	// 设定访问盐值
	salt := helper.GetClientIP(c) // crypto.UUID(16)
	ckey := "captcha-" + k + ":" + strconv.Itoa(int(info.Type)) + ":" + info.Platform + ":" + info.Account + ":" + salt
	// 验证是否发送过
	if b, err := a.Store.Check(c, ckey); err != nil {
		// 验证发送异常
		return "", err
	} else if b {
		// 发送验证码频繁， 保护后端服务器
		return "", helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-FREQUENTLY", Other: "发送频繁,稍后重试"})
	}
	var duration time.Duration = 120 * time.Second // 120秒
	a.Store.Set1(c, ckey, duration)                // 防止频繁发送

	captcha, err := info.Sender() // 发送验证码
	if err != nil {
		// 发送验证码失败
		return "", err
	}
	expired := info.Expired // 验证码有效期
	if !account.VerifySecret.Valid || account.VerifySecret.String == "" {
		// 加密密钥为空， 更新加密密钥
		account.VerifySecret.String = crypto.RandomAes32()
		account.UpdateVerifySecret(a.Sqlx)
	}
	secret := account.VerifySecret.String
	return EncryptCaptchaByAccount(c, account.ID, secret, captcha, expired)
}

// CaptchaType 解析验证类型
func (a *Signin) parseCaptchaType(c *gin.Context, b *schema.SigninOfCaptcha) (*SenderInfo, error) {
	res := &SenderInfo{
		Expired:  300 * time.Second, // 300秒, 默认验证码超时间隔
		Platform: b.Platform,        // 平台标识
	}
	if b.Mobile != "" {
		// 使用手机发送
		res.Sender = func() (string, error) {
			return a.MSender.SendCaptcha(b.Mobile)
		}
		res.Account, res.Type = b.Mobile, schema.AccountTypeMobile
	} else if b.Email != "" {
		// 使用邮箱发送
		res.Sender = func() (string, error) {
			return a.ESender.SendCaptcha(b.Email)
		}
		res.Account, res.Type = b.Email, schema.AccountTypeEmail
	} else if b.Openid != "" && b.Platform != "" {
		// 使用第三方程序发送
		res.Sender = func() (string, error) {
			return a.TSender.SendCaptcha(b.Platform, b.Openid)
		}
		res.Account, res.Type = b.Openid, schema.AccountTypeOpenid
	} else {
		// 验证码无法发送
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-NONE", Other: "验证方式无效"})
	}
	return res, nil
}

//===========================================================================================
//===========================================================================================
//===========================================================================================

// OAuth2 登陆控制
func (a *Signin) OAuth2(c *gin.Context, b *schema.SigninOfOAuth2, l func(*gin.Context, int64) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	if b.Platform != "" {
		// 当前用户
		account := schema.SigninGpaAccount{}

		// 👇👇👇 是否使用待定， 存在一定安全隐患
		defer a.saveAccountByOAuth2Code(c, b, &account)
		a.findAccountByOAuth2Code(c, b, &account)
		// 👆👆👆 是否使用待定， 存在一定安全隐患
		// 防止oauth2重复授权
		if account.ID == 0 {
			o2p := schema.OAuth2GpaPlatform{}
			if err := o2p.QueryByKID(a.Sqlx, b.Platform); err != nil {
				if sqlxc.IsNotFound(err) {
					return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NONE", Other: "无效第三方登陆"})
				}
				return nil, err
			}
			o2h, ok := a.OAuth2Selector[o2p.Type.String]
			if !ok {
				return nil, helper.NewError(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-CONTROLLER-NONE",
					Other: "无对应平台的OAuth2控制器: [{{.platform}}]"}, helper.H{
					"platform": o2p.Type,
				})
			}

			token1x := a.findUserOfToken1(c, b, &o2p)
			oauth2x := oauth2.RequestOAuth2X{
				FindHost: func() string { return o2p.SigninHost.String },
				FindUser: func(relation, openid, userid, deviceid string) (int64, error) {
					return a.findUserOfOAuth2(c, b, o2p.IsSign.Bool, o2h, &o2p, token1x, &account, relation, openid, userid, deviceid)
				},
			}
			if err := o2h.Handle(c, b, &o2p, token1x, &oauth2x); err != nil {
				if redirect, ok := err.(*helper.ErrorRedirect); ok {
					// 终止重定向， 返回json数据
					if result := c.Query("result"); result == "json" {
						status := redirect.Status
						if status == 0 {
							status = 303
						}
						// log.Println(redirect.Location)
						return nil, helper.NewSuccess(c, helper.H{
							"status":   status,
							"state":    redirect.State,
							"location": redirect.Location,
						})
					}
				}
				return nil, err
			}
		}
		if account.Status != schema.StatusEnable {
			// 账户未激活， 抛出异常
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-NOBIND", Other: "用户未绑定"})
		}

		// 获取用户信息
		b2 := &schema.SigninBody{
			Scope:    b.Scope,
			Platform: b.Platform,
			OrgCode:  b.OrgCode,
			WebToken: b.WebToken,
		}
		return a.GetSignUserInfo(c, b2, &account)
	}
	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NONE", Other: "无效第三方登陆"})
}

//===========================================================================================

// findUserOfOauth2 ...
func (a *Signin) findUserOfOAuth2(c *gin.Context, b *schema.SigninOfOAuth2, ifnew bool, o2h oauth2.Handler,
	o2p *schema.OAuth2GpaPlatform, token1x oauth2.RequestToken, account *schema.SigninGpaAccount,
	relation, openid, userid, deviceid string) (int64, error) {
	if openid == "" {
		openid = o2p.AppID.String + ":" + userid
	}
	// 查询当前登录人员身份
	if err := account.QueryByAccount(a.Sqlx, openid, schema.AccountTypeOpenid, b.Platform, b.OrgCode, false); err != nil {
		if !sqlxc.IsNotFound(err) {
			return 0, err
		}
	}
	if ifnew && account.ID == 0 {
		account.String1 = sql.NullString{Valid: true, String: deviceid} // 存储用户使用的设备信息
		var user int64
		if relation != "" && userid != "" {
			name := strings.ToUpper(relation[:1]) + ":" + userid
			account.CustomID = sql.NullString{Valid: true, String: name}
			// 通过名称或者手机号查询当前用户身份, 可以执行自动归一操作
			// if info, err := o2h.GetUser(c, o2p, token1x, relation, userid); err == nil {
			// 	log.Println(info)
			// 	a.Bus.Publish("topic:account:create", info)
			// }
		}
		if user == 0 {
			if res, _ := account.SelectByAccount(a.Sqlx, openid, schema.AccountTypeOpenid, b.Platform, "", schema.StatusEnable, 1, true); len(*res) > 0 {
				user = (*res)[0].UserID.Int64 // 用户存在归一账户， 处理归一操作， 同时允许用户登录
			}
		}
		// 用户第一次登录， 收集用户信息
		account.Account = openid
		account.AccountType = schema.AccountTypeOpenid
		account.PlatformKID = sql.NullString{Valid: b.Platform != "", String: b.Platform}
		account.OrgCode = sql.NullString{Valid: b.OrgCode != "", String: b.OrgCode}
		if user > 0 {
			account.UserID = sql.NullInt64{Valid: true, Int64: user}
			account.Status = schema.StatusEnable // 自动完成激活
		} else {
			account.Status = schema.StatusNoActivate // 待激活账户
		}
		account.UpdateAndSaveX(a.Sqlx) // 持久化， 存储账户信息
	}
	return account.ID, nil
}

func (a *Signin) findUserOfToken1(c *gin.Context, b *schema.SigninOfOAuth2, o2p *schema.OAuth2GpaPlatform) *oauth2.RequestToken1X {
	// FIXME 注意：该方法存在大数据异步调用bug， 需要同步锁处理
	// FIXME 注意：该方法存在大数据异步调用bug， 需要同步锁处理
	// FIXME 注意：该方法存在大数据异步调用bug， 需要同步锁处理
	return &oauth2.RequestToken1X{
		FindToken: func(sqlx *sqlx.DB, token *oauth2.AccessToken, platform int64) error {
			t3n := schema.OAuth2GpaAccountToken{}
			if !o2p.TokenKID.Valid {
				return nil
			} else if err := t3n.QueryByTokenKID2(a.Sqlx2, o2p.TokenKID.String); err != nil {
				return err
			}
			token.Account = t3n.AccountID.Int64
			token.Platform = o2p.ID
			token.AccessToken = t3n.AccessToken
			token.ExpiresAt = t3n.ExpiresAt
			token.RefreshToken = t3n.RefreshToken
			token.RefreshExpAt = t3n.RefreshExpAt
			token.Scope = t3n.String2
			token.AsyncLock = t3n.UpdatedAt
			return nil
		},
		SaveToken: func(sqlx *sqlx.DB, token *oauth2.AccessToken) error {
			tid := ""
			if token.Account > 0 {
				tid = jwt.NewTokenID(strconv.Itoa(int(token.Account) + 1103)) // 有用户信息令牌
				helper.SetCtxValue(c, helper.ResTknKey, tid)                  // 令牌放入缓存， 备用
			} else {
				tid = "z" + crypto.EncodeBaseX32(o2p.ID) // 无用户信息令牌
			}
			t3n := schema.OAuth2GpaAccountToken{
				TokenID:      tid,
				AccountID:    sql.NullInt64{Valid: token.Account > 0, Int64: token.Account},
				Platform:     o2p.KID,
				TokenType:    sql.NullInt32{Valid: true, Int32: 2},
				AccessToken:  token.AccessToken,
				ExpiresAt:    token.ExpiresAt,
				RefreshToken: token.RefreshToken,
				RefreshExpAt: token.RefreshExpAt,
				String2:      token.Scope,
			}
			// 对令牌进行赋值
			o3p := schema.OAuth2GpaPlatform{
				ID:       o2p.ID,
				TokenKID: sql.NullString{Valid: true, String: t3n.TokenID},
			}
			if err := t3n.UpdateAndSaveByTokenKID2(a.Sqlx2, false); err != nil {
				return err
			}
			return o3p.UpdateAndSaveByID(a.Sqlx) // 绑定令牌
		},
		LockToken: func(sqlx *sqlx.DB, platform int64) error {
			t3n := schema.OAuth2GpaAccountToken{}
			if !o2p.TokenKID.Valid {
				// 没有使用令牌
				return errors.New("no token")
			} else if err := t3n.QueryByTokenKID2(a.Sqlx2, o2p.TokenKID.String); err != nil {
				return err
			}
			t4n := schema.OAuth2GpaAccountToken{
				TokenID: t3n.TokenID,
			}
			// 更新令牌的更新时间
			return t4n.UpdateAndSaveByTokenKID2(a.Sqlx2, true)
		},
	}
}

//===========================================================================================

func (a *Signin) findAccountByOAuth2Code(c *gin.Context, b *schema.SigninOfOAuth2, account *schema.SigninGpaAccount) {
	if oac := c.Query("oac"); oac != "new" {
		val, _ := c.Cookie("zgo_oac")
		if val != "" {
			t3n := schema.OAuth2GpaAccountToken{}
			t3n.QueryByPlatformAndCode2(a.Sqlx2, b.Platform, val)
			// TokenType: 2, 标识是来自三方授权令牌
			// ErrCode: "", 标识该令牌未被主动销毁
			// CodeExpAt: > now, 标识授权有效
			if t3n.AccountID.Int64 > 0 && t3n.TokenType.Int32 == 2 && t3n.ErrCode.String == "" &&
				t3n.CodeExpAt.Valid && time.Now().Before(t3n.ExpiresAt.Time) {
				if account.QueryByID(a.Sqlx, t3n.AccountID.Int64); account.ID > 0 {
					// 直接查询上次认证信息, 该方法存在安全隐患, 但是可以减少OAuth2认证次数
					helper.SetCtxValue(c, helper.ResTknKey, val)
				}
			}
		}
	}
}

func (a *Signin) saveAccountByOAuth2Code(c *gin.Context, b *schema.SigninOfOAuth2, account *schema.SigninGpaAccount) {
	if account.ID == 0 || b.Code == "" {
		return
	}
	val, _ := c.Cookie("zgo_oac")
	tid, _ := helper.GetCtxValueToString(c, helper.ResTknKey)
	if tid != "" && val != tid {
		// 需要更新令牌， 默认授权只存放12个小时， 超时， 需要重新认证
		t3n := schema.OAuth2GpaAccountToken{
			TokenID:   tid,
			CodeToken: sql.NullString{Valid: true, String: b.Code},
			CodeExpAt: sql.NullTime{Valid: true, Time: time.Now().Add(12 * time.Hour)},
		}
		if err := t3n.UpdateAndSaveByTokenKID2(a.Sqlx2, true); err == nil {
			cke := http.Cookie{Name: "zgo_oac", Value: b.Code, Expires: time.Now().Add(12 * time.Hour)}
			http.SetCookie(c.Writer, &cke)
		}
	}
}

//===========================================================================================
