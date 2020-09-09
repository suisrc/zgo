package service

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/suisrc/zgo/modules/crypto"

	"github.com/suisrc/zgo/app/model/sqlxc"

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
	GPA                       // 数据库
	Passwd  *passwd.Validator // 密码验证其
	Store   store.Storer      // 缓存控制器
	MSender MobileSender      // 手机
	ESender EmailSender       // 邮箱
	TSender ThreeSender       // 第三方
}

// Signin 登陆控制
func (a *Signin) Signin(c *gin.Context, b *schema.SigninBody, lastSignin func(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error)) (*schema.SigninUser, error) {
	if b.Password != "" {
		return a.SigninByPasswd(c, b, lastSignin)
	}
	if b.Captcha != "" {
		if b.Code == "" {
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TYPE-CODE", Other: "校验码无效"})
		}
		return a.SigninByCaptcha(c, b, lastSignin)
	}
	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TYPE-NONE", Other: "无效登陆"})
}

// OAuth2 登陆控制
func (a *Signin) OAuth2(c *gin.Context, b *schema.SigninOfOAuth2, lastSignin func(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error)) (*schema.SigninUser, error) {
	kid := c.Param("kid")
	if kid != "" {
		b.KID = kid
		logger.Infof(c, b.KID)
	}
	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NONE", Other: "无效第三方登陆"})
}

// Captcha 发送验证码
func (a *Signin) Captcha(c *gin.Context, b *schema.SigninOfCaptcha) (string, error) {
	var duration time.Duration = 120 * time.Second // 120秒

	acc, kid, typ, expire, sender, err := a.parseCaptchaType(c, b)
	if err != nil {
		return "", err
	}
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, acc, typ, kid); err != nil {
		if sqlxc.IsNotFound(err) {
			return "", helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
		}
		return "", err
	}

	salt := helper.GetClientIP(c) // crypto.UUID(16)
	key := "captcha:" + strconv.Itoa(typ) + ":" + kid + ":" + acc + ":" + salt
	// 验证是否发送过
	if b, err := a.Store.Check(c, key); err != nil {
		return "", err
	} else if b {
		return "", helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-FREQUENTLY", Other: "发送频繁,稍后重试"})
	}
	a.Store.Set1(c, key, duration) // 防止频繁发送
	captcha, err := sender()       // 发送验证码
	if err != nil {
		return "", err
	}
	if !account.VerifySecret.Valid {
		account.VerifySecret.String = crypto.RandomAes32()
		account.UpdateVerifySecret(a.Sqlx)
	}
	return a.parseCaptchaEncrypt(c, &account, captcha, expire)
}

//============================================================================================

// SigninByPasswd 密码登陆
func (a *Signin) SigninByPasswd(c *gin.Context, b *schema.SigninBody, lastSignin func(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error)) (*schema.SigninUser, error) {
	// 查询账户信息
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, b.Username, 1, b.KID); err != nil || account.ID <= 0 {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}
	// 验证密码
	if b, err := a.VerifyPassword(b.Password, &account); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 密码验证发生异常
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	} else if !b {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}
	if b.Client != "" {
		helper.SetJwtKid(c, b.Client)
	}
	// 获取用户信息
	return a.GetSignUserByAutoRole(c, &account, b, lastSignin)
}

// SigninByCaptcha 验证码登陆
func (a *Signin) SigninByCaptcha(c *gin.Context, b *schema.SigninBody, lastSignin func(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error)) (*schema.SigninUser, error) {
	// 查询账户信息
	accountID, captchaGetter, err := a.parseCaptchaDecrypt(c, b.Code)
	if err != nil {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-CODE", Other: "校验码不正确"})
	}
	account := schema.SigninGpaAccount{}
	if err := account.QueryByID(a.Sqlx, accountID); err != nil || account.ID <= 0 {
		logger.Errorf(c, logger.ErrorWW(err)) // 账户异常
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
	} else if account.Account != b.Username || account.AccountKind.String != b.KID || !account.VerifySecret.Valid {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-USER", Other: "账户异常,请联系管理员"})
	}
	// 验证验证码
	if captcha, expire, err := captchaGetter(&account); err != nil {
		return nil, err // 解密出现问题
	} else if expire <= 0 {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-EXPIRED", Other: "验证码已过期"})
	} else if captcha != b.Captcha {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-CHECK", Other: "验证码不正确"})
	}
	if b.Client != "" {
		helper.SetJwtKid(c, b.Client)
	}
	// 获取用户信息
	return a.GetSignUserByAutoRole(c, &account, b, lastSignin)
}

//============================================================================================

// GetSignUserByAutoRole auto role
func (a *Signin) GetSignUserByAutoRole(c *gin.Context, account *schema.SigninGpaAccount, b *schema.SigninBody, lastSignin func(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error)) (*schema.SigninUser, error) {
	// 登陆用户
	suser := schema.SigninUser{}
	suser.AccountID = strconv.Itoa(account.ID)
	// 用户
	user := schema.SigninGpaUser{}
	if err := user.QueryByID(a.Sqlx, account.UserID); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户不存在"})
	} else if !user.Status {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
	}
	suser.UserName = user.Name
	suser.UserID = user.KID

	domain := b.Domain // 用户请求域名
	if domain == "" {
		domain = c.Request.Host
	}
	client := schema.JwtGpaOpts{} // 用户请求应用, 如果client.ID == ""标识client不存在
	if b.Client != "" {
		if err := client.QueryByKID(a.Sqlx, b.Client); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOEXIST", Other: "用户请求的客户端不存在"})
		}
	} else if b.Domain != "" {
		// 一般定义了重定向的域名
		if err := client.QueryByAudience(a.Sqlx, b.Domain); err != nil {
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOEXIST", Other: "用户请求的客户端不存在"})
		}
	} else {
		// do nothing
		// 不进行域名验证,但是下文会验证当前请求域名和角色域.
	}
	if client.ID > 0 {
		helper.SetJwtKid(c, client.KID) // 配置客户端
	}

	if lastSignin == nil {
		lastSignin = a.lastSigninNil
	}

	// 角色
	role := schema.SigninGpaRole{}
	if account.RoleID.Valid {
		// 单角色,账户上又绑定的角色
		if err := role.QueryByID(a.Sqlx, int(account.RoleID.Int64)); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 角色ID不存在,只有数据库数据不一致才会发生问题
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-NOROLE", Other: "用户没有有效角色[ID]"})
		}
		// 验证角色是否合法
		if !a.CheckRoleClient(c, &client, domain, &role) {
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOACCESS", Other: "用户无访问权限"})
		}
	} else if o2a, err := lastSignin(c, account.ID, client.ID); err != nil || o2a != nil && o2a.Status.Valid && o2a.Status.Bool && o2a.RoleKID.Valid {
		if err != nil {
			return nil, err
		}
		if err := role.QueryByKID(a.Sqlx, o2a.RoleKID.String); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 角色ID不存在,只有数据库数据不一致才会发生问题
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-NOROLE", Other: "用户没有有效角色[KID]"})
		}
		// 验证角色是否合法
		// if !a.CheckRoleClient(c, &client, domain, &role) {
		// 	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-NOACCESS", Other: "用户无访问权限"})
		// }
	}
	if role.ID == 0 {
		// 多角色问题, 查询第一个可用角色
		roles := []schema.SigninGpaRole{}
		if err := role.QueryByUserID(a.Sqlx, &roles, account.UserID); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		for _, r := range roles {
			if a.CheckRoleClient(c, &client, domain, &r) {
				role = r
				break
			}
		}
		if role.ID == 0 {
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
	}
	suser.RoleID = role.KID

	suser.Issuer = c.Request.Host
	suser.Audience = c.Request.Host
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

// CaptchaEncrypt 加密验证码
func (a *Signin) parseCaptchaEncrypt(c *gin.Context, account *schema.SigninGpaAccount, captcha string, expire time.Duration) (string, error) {
	checkTime := time.Now().Add(expire).Unix()                      // 过期时间
	var buffer bytes.Buffer                                         // byte buffer
	buffer.Write(crypto.Number2BytesInNetworkOrder(int(checkTime))) // 占用4个字节
	buffer.Write(crypto.Number2BytesInNetworkOrder(account.ID))     // 占用4个字节
	buffer.Write([]byte(captcha))                                   // 验证码
	// 节目加密KEY的内容
	keys, err := crypto.Base64DecodeString(account.VerifySecret.String)
	if err != nil {
		return "", err
	}
	// 对验证码进行加密, 验证码后端不存储, 加密成校验码
	checkCodeBytes, err := crypto.AesEncryptBytes(buffer.Bytes(), keys)
	if err != nil {
		return "", err // 加密出现问题
	}
	// 给出登陆账户信息,以用来进行解密
	resultCodeBytes := append(checkCodeBytes, crypto.Number2BytesInNetworkOrder(account.ID)...)
	return crypto.Base64EncodeToString(resultCodeBytes), nil
}

// CaptchaDecrypt 解密验证码
func (a *Signin) parseCaptchaDecrypt(c *gin.Context, code string) (int, func(*schema.SigninGpaAccount) (string, time.Duration, error), error) {
	resultCodeBytes, err := crypto.Base64DecodeString(code)
	if err != nil {
		return 0, nil, err
	}
	resultCodeLen := len(resultCodeBytes)
	if resultCodeLen < 12 {
		// 不可预知异常, 往往来自恶意攻击
		return 0, nil, errors.New("code is error")
	}
	accountID := crypto.BytesNetworkOrder2Number(resultCodeBytes[resultCodeLen-4:])
	return accountID, func(account *schema.SigninGpaAccount) (string, time.Duration, error) {
		keys, err := crypto.Base64DecodeString(account.VerifySecret.String)
		if err != nil {
			return "", 0, err
		}
		checkCodeBytes, err := crypto.AesDecryptBytes(resultCodeBytes[:resultCodeLen-4], keys)
		accountID2 := crypto.BytesNetworkOrder2Number(checkCodeBytes[4:8])
		if accountID != accountID2 {
			// 不可预知异常, 往往来自恶意攻击
			return "", 0, errors.New("account id error")
		}
		checkTime := crypto.BytesNetworkOrder2Number(checkCodeBytes[:4])
		expire := time.Unix(int64(checkTime), 0).Sub(time.Now())
		captcha := string(checkCodeBytes[8:])
		return captcha, expire, nil
	}, nil
}

// CaptchaType 解析验证类型
func (a *Signin) parseCaptchaType(c *gin.Context, b *schema.SigninOfCaptcha) (string, string, int, time.Duration, func() (string, error), error) {
	var expire time.Duration = 300 * time.Second // 300秒
	var acc string
	var typ int
	var sender func() (string, error)
	if b.Mobile != "" {
		// 使用手机发送
		sender = func() (string, error) {
			return a.MSender.SendCaptcha(b.Mobile)
		}
		acc, typ = b.Mobile, int(schema.ATMobile)
	} else if b.Email != "" {
		// 使用邮箱发送
		sender = func() (string, error) {
			return a.ESender.SendCaptcha(b.Email)
		}
		acc, typ = b.Email, int(schema.ATEmail)
	} else if b.Openid != "" && b.KID != "" {
		// 使用第三方程序发送
		sender = func() (string, error) {
			return a.TSender.SendCaptcha(b.KID, b.Openid)
		}
		acc, typ = b.Openid, int(schema.ATOpenid)
	} else {
		// 验证码无法发送
		return "", "", 0, 0, nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CAPTCHA-NONE", Other: "验证方式无效"})
	}
	return acc, b.KID, typ, expire, sender, nil
}

//============================================================================================

// GetSignUserBySelectRole 通过角色选择获取用户信息
// 用户选择角色,不对角色进行如何判定,当发现用户具有多角色的时候,用户选择使用的角色
// 这种方式适合单系统,由于多系统涉及到用于域的概念,是很难完成多角色的自由切换
func (a *Signin) GetSignUserBySelectRole(c *gin.Context, account *schema.SigninGpaAccount, b *schema.SigninBody) (*schema.SigninUser, error) {
	// 登陆用户
	suser := schema.SigninUser{}
	suser.AccountID = strconv.Itoa(account.ID)
	// 用户
	user := schema.SigninGpaUser{}
	if err := user.QueryByID(a.Sqlx, account.UserID); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户不存在"})
	} else if !user.Status {
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
	suser.RoleID = role.KID

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

// 空的上次登陆验证器
func (a *Signin) lastSigninNil(c *gin.Context, aid, cid int) (*schema.SigninGpaOAuth2Account, error) {
	return nil, nil
}
