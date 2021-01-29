package service

import (
	"regexp"
	"strconv"
	"time"

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
}

// Signin 登陆控制
// params: c 访问上下文
// params: b 请求参数
// params: l 验证最后一次登录结果
// result: 登录者信息
func (a *Signin) Signin(c *gin.Context, b *schema.SigninBody, l func(*gin.Context, int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
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
	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TYPE-NONE", Other: "无效登陆"})
}

// OAuth2 登陆控制
// func (a *Signin) OAuth2(c *gin.Context, b *schema.SigninOfOAuth2, last func(c *gin.Context, aid, cid int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
// 	if b.KID == "" {
// 		b.KID = c.Param("kid")
// 	}
// 	if b.KID != "" {
// 		o2p := schema.SigninGpaOAuth2Platfrm{}
// 		if err := o2p.QueryByKID(a.Sqlx, b.KID); err != nil {
// 			return nil, err
// 		}
// 		o2h, ok := a.OAuth2Selector[o2p.Platform]
// 		if !ok {
// 			return nil, helper.NewError(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-NOP6M",
// 				Other: "无有效登陆控制器: [{{.platfrom}}]"}, helper.H{
// 				"platfrom": o2p.Platform,
// 			})
// 		}
//
// 		// 当前用户
// 		account := &schema.SigninGpaAccount{}
// 		if err := o2h.Handle(c, b, &o2p, true, func(openid string) (int, error) {
// 			// 查询当前登录人员身份
// 			if err := account.QueryByAccount(a.Sqlx, openid, int(schema.ATOpenid), o2p.KID); err != nil {
// 				if sqlxc.IsNotFound(err) {
// 					return 0, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-NOBIND", Other: "用户未绑定"})
// 				}
// 				return 0, err
// 			}
// 			if account.ID == 0 {
// 				return 0, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-NOBIND", Other: "用户未绑定"})
// 			}
// 			return account.ID, nil
// 		}); err != nil {
// 			if redirect, ok := err.(*helper.ErrorRedirect); ok {
// 				if result := c.Query("result"); result == "json" {
// 					code := redirect.Code
// 					if code == 0 {
// 						code = 303
// 					}
// 					return nil, helper.NewSuccess(c, helper.H{
// 						"code":     code,
// 						"location": redirect.Location,
// 					})
// 				}
// 			}
// 			return nil, err
// 		}
// 		// 获取用户信息
// 		return a.GetSignUserByAutoRole(c, account, b.Domain, b.Client, last)
// 	}
// 	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NONE", Other: "无效第三方登陆"})
// }

//============================================================================================

// SigninByPasswd 密码登陆
func (a *Signin) SigninByPasswd(c *gin.Context, b *schema.SigninBody, l func(*gin.Context, int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	// 查询账户信息
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, b.Username, int(schema.AccountTypeName), b.KID, b.Org); err != nil || account.ID <= 0 {
		// 无法查询到账户， 是否可以使用 2(手机)， 3(邮箱) 查询， 待定
		account.ID = 0
		if p, _ := regexp.MatchString(`^(1[3-8]\d{9}$`, b.Username); p {
			// 使用手机方式登录(只匹配中国手机号)
			err = account.QueryByParentAccount(a.Sqlx, b.Username, int(schema.AccountTypeMobile), b.KID, b.Org)
		} else if p, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+$`, b.Username); p {
			// 使用邮箱方式登录
			err = account.QueryByParentAccount(a.Sqlx, b.Username, int(schema.AccountTypeEmail), b.KID, b.Org)
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
	// 获取用户信息
	return a.GetSignUserInfo(c, &account, l)
}

// SigninByCaptcha 验证码登陆
func (a *Signin) SigninByCaptcha(c *gin.Context, b *schema.SigninBody, l func(*gin.Context, int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
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
	} else if account.Account != b.Username || account.PlatformKID.String != b.KID || !account.VerifySecret.Valid {
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
	// 获取用户信息
	return a.GetSignUserInfo(c, &account, l)
}

//============================================================================================

// GetSignUserInfo with role
func (a *Signin) GetSignUserInfo(c *gin.Context, sa *schema.SigninGpaAccount, l func(*gin.Context, int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	if sa.Status != schema.StatusEnable { // 账户被禁用
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ACCOUNT-DISABLE", Other: "账户被禁用,请联系管理员"})
	}
	if l == nil {
		l = a.lastSignInNil
	}
	// 登陆用户
	suser := schema.SigninUser{}
	//suser.AccountID = strconv.Itoa(sa.ID) // SigninUser -> 1
	//suser.UserIdxID = strconv.Itoa(sa.UserID)
	suser.TokenID, _ = helper.GetCtxValueToString(c, helper.ResTokenKey) //  配置系统给定的TokenID
	if suser.TokenID == "" {                                             // account加密需要令牌， 所以令牌不能为空
		suser.TokenID = jwt.NewTokenID(strconv.Itoa(sa.ID + 1103))
	}
	suser.Account, _ = EncryptAccountWithUser(c, sa.ID, int(sa.UserID.Int64), suser.TokenID) // 账户信息
	if err := a.SetSignUserWithUser(c, sa, &suser); err != nil {                             // 用户信息
		return nil, err
	}
	if err := a.SetSignUserWithClient(c, sa, &suser); err != nil { // 访问令牌签名
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
// 如果一个人具有管理员权限， 其所有的角色都会被舍弃， 只保留管理员角色
func (a *Signin) SetSignUserWithRole(c *gin.Context, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	// 查询用户的所有的角色
	if !sa.UserID.Valid {
		if sa.RoleID.Valid {
			grr := schema.SigninGpaRole{}
			// 登录账户上绑定了角色
			if err := grr.QueryByRoleAndOrg(a.Sqlx, int(sa.RoleID.Int64), suser.OrgCode); err != nil {
				if sqlxc.IsNotFound(err) { // 角色没有找到
					return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-NOROLE", Other: "账户角色失效，账户无法使用"})
				}
				return err // 未知异常
			}
			if grr.OrgAdm {
				suser.OrgAdmin = schema.SuperUser // 超级管理员， 不需要角色
			} else if grr.SvcCode.Valid {
				suser.SetUserRoles([]string{grr.SvcCode.String + ":" + grr.Name}) // 应用角色
			} else {
				suser.SetUserRoles([]string{grr.Name}) // 租户角色
			}
		}
		return nil
	} else if sa.RoleID.Valid {
		gur := schema.SigninGpaUserRole{}
		// 登录账户上绑定了角色
		if err := gur.QueryByUserAndRoleAndOrg(a.Sqlx, int(sa.UserID.Int64), int(sa.RoleID.Int64), suser.OrgCode); err != nil {
			if sqlxc.IsNotFound(err) { // 角色没有找到
				return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-NOROLE", Other: "账户角色失效，账户无法使用"})
			}
			return err // 未知异常
		}
		if gur.OrgAdm {
			suser.OrgAdmin = schema.SuperUser
			return nil // 超级管理员， 不需要角色
		} else if gur.SvcCode.Valid {
			// 应用角色
			suser.SetUserRoles([]string{gur.SvcCode.String + ":" + gur.Name})
		} else {
			// 租户角色
			suser.SetUserRoles([]string{gur.Name})
		}
		return nil
	}
	// 账户上没有角色， 取用户在对应租户下的所有角色
	gur := schema.SigninGpaUserRole{}
	if roles, err := gur.QueryAllByUserAndOrg(a.Sqlx, int(sa.UserID.Int64), suser.OrgCode); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
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
				// 应用角色
				rs = append(rs, r.SvcCode.String+":"+r.Name)
			} else {
				// 租户角色
				rs = append(rs, r.Name)
			}
		}
		// 设定用户角色
		suser.SetUserRoles(rs)
	}
	return nil
}

// SetSignUserWithClient with client info
func (a *Signin) SetSignUserWithClient(c *gin.Context, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	// TODO JWT多加密配置方案
	// domain := c.Request.Host // 用户请求域名
	// client := schema.JwtGpaOpts{}
	// if err := client.QueryByAudience(a.Sqlx, domain, suser.OrgCode); err == nil && client.ID > 0 {
	// 	// return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOEXIST", Other: "用户请求的客户端不存在"})
	// 	helper.SetCtxValue(c, helper.ResJwtKey, client.KID) // 配置客户端, 该内容会影响JWT签名方式
	// 	if client.Issuer.Valid {
	// 		suser.Issuer = client.Issuer.String
	// 	}
	// 	if client.Audience.Valid {
	// 		suser.Audience = client.Audience.String
	// 	}
	// }
	if suser.Issuer == "" {
		suser.Issuer = c.Request.Host
	}
	if suser.Audience == "" {
		suser.Audience = c.Request.Host
	}
	return nil
}

// SetSignUserWithUser with user info
func (a *Signin) SetSignUserWithUser(c *gin.Context, sa *schema.SigninGpaAccount, suser *schema.SigninUser) error {
	if !sa.UserID.Valid {
		// 账户上没有用户信息， 待验证账户， 允许登录
		suser.OrgCode = sa.OrgCod.String
		return nil
	}
	if sa.OrgCod.Valid {
		// 账户上绑定了租户， 使用用户的租户账户
		user := schema.SigninGpaOrgUser{}
		if err := user.QueryByUserAndOrg(a.Sqlx, int(sa.UserID.Int64), sa.OrgCod.String); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户信息发生异常"})
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
		// 使用用户的平台账户
		user := schema.SigninGpaUser{}
		if err := user.QueryByID(a.Sqlx, int(sa.UserID.Int64), ""); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户信息发生异常"})
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
	if err := account.QueryByAccount(a.Sqlx, info.Acc, info.Typ, info.KID, b.Org); err != nil {
		if sqlxc.IsNotFound(err) {
			return "", helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
		}
		return "", err
	}
	// 设定访问盐值
	salt := helper.GetClientIP(c) // crypto.UUID(16)
	ckey := "captcha-" + k + ":" + strconv.Itoa(info.Typ) + ":" + info.KID + ":" + info.Acc + ":" + salt
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
	expire := info.Expire // 验证码有效期
	if !account.VerifySecret.Valid || account.VerifySecret.String == "" {
		// 加密密钥为空， 更新加密密钥
		account.VerifySecret.String = crypto.RandomAes32()
		account.UpdateVerifySecret(a.Sqlx)
	}
	secret := account.VerifySecret.String
	return EncryptCaptchaByAccount(c, account.ID, secret, captcha, expire)
}

// CaptchaType 解析验证类型
func (a *Signin) parseCaptchaType(c *gin.Context, b *schema.SigninOfCaptcha) (*SenderInfo, error) {
	res := &SenderInfo{
		Expire: 300 * time.Second, // 300秒, 默认验证码超时间隔
		KID:    b.KID,             // 平台标识
	}
	if b.Mobile != "" {
		// 使用手机发送
		res.Sender = func() (string, error) {
			return a.MSender.SendCaptcha(b.Mobile)
		}
		res.Acc, res.Typ = b.Mobile, int(schema.AccountTypeMobile)
	} else if b.Email != "" {
		// 使用邮箱发送
		res.Sender = func() (string, error) {
			return a.ESender.SendCaptcha(b.Email)
		}
		res.Acc, res.Typ = b.Email, int(schema.AccountTypeEmail)
	} else if b.Openid != "" && b.KID != "" {
		// 使用第三方程序发送
		res.Sender = func() (string, error) {
			return a.TSender.SendCaptcha(b.KID, b.Openid)
		}
		res.Acc, res.Typ = b.Openid, int(schema.AccountTypeOpenid)
	} else {
		// 验证码无法发送
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-NONE", Other: "验证方式无效"})
	}
	return res, nil
}

//============================================================================================

// 空的上次登陆验证器
func (a *Signin) lastSignInNil(c *gin.Context, aid int) (*schema.SigninGpaAccountToken, error) {
	return nil, nil
}
