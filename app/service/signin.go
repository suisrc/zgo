package service

import (
	"strconv"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	gi18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/passwd"

	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
)

// Signin 账户管理
type Signin struct {
	GPA                     // 数据库
	Passwd passwd.Validator // 密码验证其
}

//============================================================================================

// SigninByPasswd 密码登陆
func (a *Signin) SigninByPasswd(c *gin.Context, b *schema.SigninBody) (*schema.SigninUser, error) {
	// 查询账户信息
	account := schema.SigninGpaAccount{}
	if err := account.QueryByAccount(a.Sqlx, b.Username, 1, b.KID); err != nil {
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
	return a.GetSignUserBySelectRole(c, &account, b)
}

// GetSignUserBySelectRole 通过角色选择获取用户信息
// 用户选择角色,不对角色进行如何判定,当发现用户具有多角色的时候,用户选择使用的角色
// 这种方式适合单系统,由于多系统涉及到用于域的概念,是很难完成多角色的自由切换
func (a *Signin) GetSignUserBySelectRole(c *gin.Context, account *schema.SigninGpaAccount, b *schema.SigninBody) (*schema.SigninUser, error) {
	// 登陆用户
	suser := schema.SigninUser{}
	if account.ID > 0 {
		suser.AccountID = strconv.Itoa(account.ID)
	}
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
	if account.RoleID.Valid {
		// 账户已经绑定角色, 失败账户角色登陆
		// 注意Role和User是通过UserRole关联,除此之外,可以通过直接配置账户,以使他具有固定的角色,该绑定关系脱离user_role管理
		role := schema.SigninGpaRole{}
		if err := role.QueryByID(a.Sqlx, int(account.RoleID.Int64)); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色[ID]"})
		}
		suser.RoleID = role.KID
	} else if b.Role != "" {
		role := schema.SigninGpaRole{}
		if err := role.QueryByKID(a.Sqlx, b.Role); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色[KID]"})
		}
		suser.RoleID = role.KID
	} else {
		// 多角色问题
		role := schema.SigninGpaRole{}
		roles := []schema.SigninGpaRole{}
		if err := role.QueryByUserID(a.Sqlx, &roles, account.UserID); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色[MUL]"})
		}
		switch len(roles) {
		case 0:
			// 没有角色,赋值默认角色
			// do nothings, 目前默认角色问题已经迁移到[norole]问题中处理
		case 1:
			suser.RoleID = roles[0].KID
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

	suser.Issuer = c.Request.Host
	suser.Audience = c.Request.Host
	return &suser, nil
}

// VerifyPassword 验证密码
//============================================================================================
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
