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
	gi18n "github.com/suisrc/gin-i18n"
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
	if err := account.QueryByAccount(a.Sqlx, b.Username, int(schema.AccountTypeName), b.KID); err != nil || account.ID <= 0 {
		// 无法查询到账户， 是否可以使用 2(手机)， 3(邮箱) 查询， 待定
		account.ID = 0
		if p, _ := regexp.MatchString(`^(1[3-8]\d{9}$`, b.Username); p {
			// 使用手机方式登录(只匹配中国手机号)
			err = account.QueryByParentAccount(a.Sqlx, b.Username, int(schema.AccountTypeMobile), b.KID)
		} else if p, _ := regexp.MatchString(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+$`, b.Username); p {
			// 使用邮箱方式登录
			err = account.QueryByParentAccount(a.Sqlx, b.Username, int(schema.AccountTypeEmail), b.KID)
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
	return a.GetSignUserWithRole(c, &account, l)
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
	return a.GetSignUserWithRole(c, &account, l)
}

//============================================================================================

// GetSignUserWithRole with role
func (a *Signin) GetSignUserWithRole(c *gin.Context, sa *schema.SigninGpaAccount, l func(*gin.Context, int) (*schema.SigninGpaAccountToken, error)) (*schema.SigninUser, error) {
	if sa.Status != schema.StatusEnable {
		// 账户被禁用
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
	if suser.TokenID == "" {
		// account加密需要令牌， 所以令牌不能为空
		suser.TokenID = jwt.NewTokenID(strconv.Itoa(sa.ID + 1103))
	}
	var err error
	suser.Account, err = EncryptAccountWithUser(c, sa.ID, sa.UserID, suser.TokenID)
	if err != nil {
		return nil, err
	}

	if sa.OrgCod.Valid {
		// 账户上绑定了租户， 使用用户的租户账户
		user := schema.SigninGpaOrgUser{}
		if err := user.QueryByUserAndOrg(a.Sqlx, sa.UserID, sa.OrgCod.String); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户信息发生异常"})
		} else if user.Status != schema.StatusEnable {
			// 租户账户被禁用
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
		}
		suser.UserID = user.UnionKID
		suser.UserName = user.Name
		suser.OrgCode = user.OrgCode
		suser.OrgUsrID = user.CustomID.String
		// suser.OrgAdmin = "admin" // 用户角色为指定
	} else {
		// 使用用户的平台账户
		user := schema.SigninGpaUser{}
		if err := user.QueryByID(a.Sqlx, sa.UserID, ""); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户信息发生异常"})
		} else if user.Status != schema.StatusEnable {
			// 平台账户被禁用
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
		}
		suser.UserID = user.KID
		suser.UserName = user.Name
	}
	//
	//domain := bDomain // 用户请求域名
	//if domain == "" {
	//	domain = c.Request.Host
	//}
	//client := schema.JwtGpaOpts{} // 用户请求应用, 如果client.ID == ""标识client不存在
	//if bClient != "" {
	//	if err := client.QueryByKID(a.Sqlx, bClient); err != nil {
	//		logger.Errorf(c, logger.ErrorWW(err))
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOEXIST", Other: "用户请求的客户端不存在"})
	//	}
	//} else if bDomain != "" {
	//	// 一般定义了重定向的域名
	//	if err := client.QueryByAudience(a.Sqlx, bDomain); err != nil {
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOEXIST", Other: "用户请求的客户端不存在"})
	//	}
	//} else {
	//	// do nothing
	//	// 不进行域名验证,但是下文会验证当前请求域名和角色域.
	//}
	//if client.ID > 0 {
	//	helper.SetCtxValue(c, helper.ResJwtKey, client.KID) // 配置客户端, 该内容会影响JWT签名方式
	//	if client.Issuer.Valid {
	//		suser.Issuer = client.Issuer.String // SigninUser -> 4
	//	}
	//	if client.Audience.Valid {
	//		suser.Audience = client.Audience.String // SigninUser -> 5
	//	}
	//}
	if suser.Issuer == "" {
		suser.Issuer = c.Request.Host
	}
	if suser.Audience == "" {
		suser.Audience = c.Request.Host
	}

	//// 角色
	//role := schema.SigninGpaRole{}
	//if account.RoleID.Valid {
	//	// 单角色,账户上又绑定的角色
	//	if err := role.QueryByID(a.Sqlx, int(account.RoleID.Int64)); err != nil {
	//		logger.Errorf(c, logger.ErrorWW(err)) // 角色ID不存在,只有数据库数据不一致才会发生问题
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-NOROLE", Other: "用户没有有效角色[ID]"})
	//	}
	//	// 验证角色是否合法
	//	if !a.CheckRoleClient(c, &client, domain, &role) {
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOACCESS", Other: "用户无访问权限"})
	//	}
	//} else if o2a, err := last(c, account.ID, client.ID); err != nil || o2a != nil && o2a.Status.Valid && o2a.Status.Bool && o2a.RoleKID.Valid {
	//	if err != nil {
	//		return nil, err
	//	}
	//	if err := role.QueryByKID(a.Sqlx, o2a.RoleKID.String); err != nil {
	//		logger.Errorf(c, logger.ErrorWW(err)) // 角色ID不存在,只有数据库数据不一致才会发生问题
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-NOROLE", Other: "用户没有有效角色[KID]"})
	//	}
	//	// 验证角色是否合法
	//	// if !a.CheckRoleClient(c, &client, domain, &role) {
	//	// 	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOACCESS", Other: "用户无访问权限"})
	//	// }
	//}
	//if role.ID == 0 {
	//	// 多角色问题, 查询第一个可用角色
	//	roles := []schema.SigninGpaRole{}
	//	if err := role.QueryByUserID(a.Sqlx, &roles, account.UserID); err != nil {
	//		logger.Errorf(c, logger.ErrorWW(err))
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
	//	}
	//	for _, r := range roles {
	//		if a.CheckRoleClient(c, &client, domain, &r) {
	//			role = r
	//			break
	//		}
	//	}
	//	if role.ID == 0 {
	//		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
	//	}
	//}
	//suser.Role = role.KID                                               // 访问令牌KID
	return &suser, nil
}

//============================================================================================

// CheckRoleClient 确定访问权限
func (a *Signin) CheckRoleClient(c *gin.Context, client *schema.JwtGpaOpts, domain string, role *schema.SigninGpaRole) bool {
	if role.ID == 0 {
		return false // 没有可用角色
	}
	if client.ID == 0 {
		// 未指定客户端信息, 使用的角色域必须为空或者等于请求域
		return role.Domain == nil || *role.Domain == domain
	}
	if client.Audience.Valid && client.Audience.String != domain {
		return false // 客户端, 接受域必须等于请求域
	}
	if role.Domain == nil {
		return true // 该角色为全平台角色
	}
	if client.Audience.Valid && client.Audience.String != *role.Domain {
		return false // 该角色的作用域和客户端的作用域冲突
	}
	// c.Redirect
	return *role.Domain == domain // 该角色的域和请求域必须相等
}

//============================================================================================

// GetSignUserBySelectRole 通过角色选择获取用户信息
// 用户选择角色,不对角色进行如何判定,当发现用户具有多角色的时候,用户选择使用的角色
// 这种方式适合单系统,由于多系统涉及到用于域的概念,是很难完成多角色的自由切换
func (a *Signin) GetSignUserBySelectRole(c *gin.Context, account *schema.SigninGpaAccount, b *schema.SigninBody) (*schema.SigninUser, error) {
	// 登陆用户
	suser := schema.SigninUser{}
	suser.Account = strconv.Itoa(account.ID)
	// 用户
	user := schema.SigninGpaUser{}
	if err := user.QueryByID(a.Sqlx, account.UserID, ""); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户不存在"})
	} else if user.Status != 1 {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
	}
	suser.UserName = user.Name
	suser.UserID = user.KID

	// 角色
	role := schema.SigninGpaRole{}
	if account.RoleID.Valid {
		// 账户已经绑定角色, 失败账户角色登陆
		// 注意Role和User是通过UserRole关联,除此之外,可以通过直接配置账户,以使他具有固定的角色,该绑定关系脱离user_role管理
		if err := role.QueryByID(a.Sqlx, int(account.RoleID.Int64)); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色[ID]"})
		}
	} else if b.Role != "" {
		role := schema.SigninGpaRole{}
		if err := role.QueryByKID(a.Sqlx, b.Role); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色[KID]"})
		}
	} else {
		// 多角色问题
		roles := []schema.SigninGpaRole{}
		if err := role.QueryByUserID(a.Sqlx, &roles, account.UserID); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		switch len(roles) {
		case 0:
			// 没有角色,赋值默认角色
			// do nothings, 目前默认角色问题已经迁移到[norole]问题中处理
		case 1:
			role = roles[0]
		default:
			// 用户有多角色
			// return nil, helper.NewError(c, helper.ShowWarn, "WARN-SIGNIN-ROLE-MULTI-ERROR", "多角色")
			return nil, helper.NewSuccess(c, map[string]interface{}{
				"status":  "error", // 登陆失败
				"message": gi18n.FormatText(c, &i18n.Message{ID: "service.signin.select-role-text", Other: "请选择角色"}),
				"roles":   roles,
			})
		}
	}
	// suser.UserRole = role.KID

	suser.Issuer = c.Request.Host
	suser.Audience = c.Request.Host
	return &suser, nil
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
	if err := account.QueryByAccount(a.Sqlx, info.Acc, info.Typ, info.KID); err != nil {
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
