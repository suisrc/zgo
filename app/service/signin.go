package service

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/modules/helper"

	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
)

// Signin 账户管理
type Signin struct {
	GPA
}

// Signin 登入
//============================================================================================
func (a *Signin) Signin(c *gin.Context, b *schema.SigninBody) (*schema.SigninUser, error) {

	// 查询账户信息
	account := schema.AccountSignin{}
	err := a.GPA.Sqlx.Get(&account, `select id, verify_type, password, password_salt, password_type, user_id, role_id
		from account where account=? and account_type='user' and platform='ZGO' and status=1`, b.Username)
	if err != nil {
		// logger.Errorf(c, err.Error()) // 未找对应的用户
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}
	// 验证密码
	if b, err := a.verifyPassword(b.Password, &account); err != nil || !b {
		// logger.Errorf(c, err.Error()) // 密码出现问题
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-PASSWD-ERROR", Other: "用户或密码错误"})
	}

	suser := schema.SigninUser{}
	// 客户端
	if b.Client != "" {
		client := schema.ClientSignin{}
		err := a.GPA.Sqlx.Get(&client, "select id, issuer, audience from user where client_key=?", b.Client)
		if err != nil || !client.Issuer.Valid || !client.Audience.Valid {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-CLIENT-ERROR", Other: "客户端错误"})
		}

		suser.Issuer = client.Issuer.String
		suser.Audience = client.Audience.String
	} else {
		suser.Issuer = c.Request.Host
		suser.Audience = c.Request.Host
	}

	// 角色
	if account.RoleID.Valid {
		role := schema.RoleSignin{}
		err = a.GPA.Sqlx.Get(&role, "select id, uid, name from role where id=?", account.RoleID)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		suser.RoleID = role.UID
	} else if b.Role != "" {
		role := schema.RoleSignin{}
		err = a.GPA.Sqlx.Get(&role, "select id, uid, name from role where uid=?", b.Role)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		suser.RoleID = role.UID
	} else {
		// 多角色问题
		roles := []schema.RoleSignin{}
		err = a.GPA.Sqlx.Select(&roles, `select r.id, r.uid, r.name from user_role ur inner join role r on r.id=ur.role_id where ur.user_id=?`, account.UserID)
		if err != nil {
			logger.Errorf(c, err.Error())
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
		}
		switch len(roles) {
		case 0:
			// 没有角色,赋值默认角色
			role := schema.RoleSignin{}
			err = a.GPA.Sqlx.Get(&role, "select id, uid, name from role where name=?", "default")
			if err != nil {
				logger.Errorf(c, err.Error())
				return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-ROLE-ERROR", Other: "用户没有有效角色"})
			}
			suser.RoleID = role.UID
		case 1:
			suser.RoleID = roles[0].UID
		default:
			// 用户有多角色
			// return nil, helper.NewError(c, helper.ShowWarn, "WARN-SIGNIN-ROLE-MULTI-ERROR", "多角色")
			return nil, helper.NewSuccess(c, map[string]interface{}{
				"status":  "error", // 登陆失败
				"message": helper.FormatText(c, &i18n.Message{ID: "service.signin.select-role-text", Other: "请选择角色"}),
				"roles":   roles,
			})
		}
	}

	// 用户
	user := schema.UserSignin{}
	err = a.GPA.Sqlx.Get(&user, "select id, uid, name from user where id=?", account.UserID)
	if err != nil {
		logger.Errorf(c, err.Error())
		return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-USER-ERROR", Other: "用户不存在"})
	}
	suser.UserName = user.Name
	suser.UserID = user.UID

	return &suser, nil
}

// 验证密码
//============================================================================================
func (a *Signin) verifyPassword(pwd string, acc *schema.AccountSignin) (bool, error) {
	result := pwd == acc.Password.String
	return result, nil
}
