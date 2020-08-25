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

// Signin 登入
//============================================================================================
func (a *Signin) Signin(c *gin.Context, b *schema.SigninBody) (*schema.SigninUser, error) {

	// 查询账户信息
	account := schema.SigninGpaAccount{}
	err := a.GPA.Sqlx.Get(&account, account.SQLByAccount(), b.Username)
	if err != nil {
		// logger.Errorf(c, err.Error()) // 未找对应的用户
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}
	// 验证密码
	if b, err := a.verifyPassword(b.Password, &account); err != nil {
		logger.Errorf(c, err.Error()) // 密码验证发生异常
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	} else if !b {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}

	suser := schema.SigninUser{}
	if account.ID > 0 {
		suser.SIID = strconv.Itoa(account.ID)
	}
	// 用户
	user := schema.SigninGpaUser{}
	err = a.GPA.Sqlx.Get(&user, user.SQLByID(), account.UserID)
	if err != nil {
		logger.Errorf(c, err.Error()) // 这里发生不可预知异常,登陆账户存在,但是账户对用的用户不存在
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户不存在"})
	} else if !user.Status {
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-DISABLE", Other: "用户被禁用,请联系管理员"})
	}
	suser.UserName = user.Name
	suser.UserID = user.UID

	// 角色
	if account.RoleID.Valid {
		role := schema.SigninGpaRole{}
		err = a.GPA.Sqlx.Get(&role, role.SQLByID(), account.RoleID)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		suser.RoleID = role.UID
	} else if b.Role != "" {
		role := schema.SigninGpaRole{}
		err = a.GPA.Sqlx.Get(&role, role.SQLByUID(), b.Role)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		suser.RoleID = role.UID
	} else {
		// 多角色问题
		role := schema.SigninGpaRole{}
		roles := []schema.SigninGpaRole{}
		err = a.GPA.Sqlx.Select(&roles, role.SQLByUserID(), account.UserID)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		switch len(roles) {
		case 0:
			// 没有角色,赋值默认角色
			// do nothings, 目前默认角色问题已经迁移到[norole]问题中处理
			// role := schema.SigninGpaRole{}
			// err = a.GPA.Sqlx.Get(&role, role.SQLByName(), "default")
			// if err != nil {
			// 	logger.Errorf(c, err.Error())
			// 	return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
			// }
			// suser.RoleID = role.UID
		case 1:
			suser.RoleID = roles[0].UID
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

// 验证密码
//============================================================================================
func (a *Signin) verifyPassword(pwd string, acc *schema.SigninGpaAccount) (bool, error) {
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

// Left 输入的密码
func (a *PasswdCheck) Left() string {
	return a.Password
}

// Right 保存的加密密码
func (a *PasswdCheck) Right() string {
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
