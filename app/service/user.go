package service

import (
	"github.com/suisrc/zgo/modules/store"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/oauth2"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/helper"
)

// User 账户管理
type User struct {
	gpa.GPA                        // 数据库
	Store          store.Storer    // 缓存控制器
	OAuth2Selector oauth2.Selector // OAuth2选择器
}

// Bind 绑定第三方账户
func (a *User) Bind(c *gin.Context, b *schema.SigninOfOAuth2) error {
	// suser, sok := helper.GetUserInfo(c)
	// if !sok {
	// 	helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-BIND-NOUSER", Other: "无法获取当前有效用户"})
	// }
	//
	// if b.KID == "" {
	// 	b.KID = c.Param("kid")
	// }
	// if b.KID != "" {
	// 	o2p := schema.SigninGpaOAuth2Platfrm{}
	// 	if err := o2p.QueryByKID(a.Sqlx, b.KID); err != nil {
	// 		return err
	// 	}
	// 	o2h, ok := a.OAuth2Selector[o2p.Platform]
	// 	if !ok {
	// 		return helper.NewError(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-NOP6M",
	// 			Other: "无有效登陆控制器: [{{.platfrom}}]"}, helper.H{
	// 			"platfrom": o2p.Platform,
	// 		})
	// 	}
	//
	// 	if err := o2h.Handle(c, b, &o2p, true, func(openid string) (int, error) {
	// 		user := schema.SigninGpaUser{}
	// 		if err := user.QueryByKID(a.Sqlx, suser.GetUserID()); err != nil {
	// 			return 0, err
	// 		}
	// 		// 查询账户是否被使用
	// 		account := &schema.SigninGpaAccount{}
	// 		if err := account.QueryByAccountSkipStatus(a.Sqlx, openid, int(schema.ATOpenid), o2p.KID); err == nil {
	// 			if user.ID == account.UserID {
	// 				return 0, nil // 已经完成了绑定
	// 			}
	// 			return 0, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-BOUND", Other: "当前用户已经绑定其他账户"})
	// 		} else if !sqlxc.IsNotFound(err) {
	// 			return 0, err // 数据库发生异常
	// 		} else if account.ID > 0 {
	// 			return 0, errors.New("unkown error") // 未知异常
	// 		}
	// 		account.Account = openid
	// 		account.AccountType = int(schema.ATOpenid)
	// 		account.AccountKind = sql.NullString{Valid: true, String: o2p.KID}
	// 		account.UserID = user.ID
	//
	// 		if err := account.UpdateAndSaveX(a.Sqlx); err != nil { // 增加绑定的登陆凭据
	// 			return 0, err
	// 		}
	// 		// 结束
	// 		return account.ID, nil
	// 	}); err != nil {
	// 		if redirect, ok := err.(*helper.ErrorRedirect); ok {
	// 			if token, err := jwt.GetBearerToken(c); err != nil {
	// 				return err
	// 			} else if token != "" {
	// 				// 重定向回调会丢失我们的令牌,可以使用state来冲缓存中还原我们的令牌,需要配合wire_auth中的tokenFunc方法使用
	// 				a.Store.Set(c, redirect.State, token, time.Duration(300)*time.Second)
	// 			}
	// 			if result := c.Query("result"); result == "json" {
	// 				code := redirect.Code
	// 				if code == 0 {
	// 					code = 303
	// 				}
	// 				return helper.NewSuccess(c, helper.H{
	// 					"code":     code,
	// 					"location": redirect.Location,
	// 				})
	// 			}
	// 		}
	// 		return err
	// 	}
	// 	return nil
	// }

	return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NONE", Other: "无效第三方登陆"})
}

// Unbind 解绑第三方账户
func (a *User) Unbind(c *gin.Context, b *schema.SigninOfOAuth2) error {
	// suser, sok := helper.GetUserInfo(c)
	// if !sok {
	// 	helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-UNBIND-NOUSER", Other: "无法获取当前有效用户"})
	// }
	// if b.KID == "" {
	// 	b.KID = c.Param("kid")
	// }
	// if b.KID != "" {
	// 	user := schema.SigninGpaUser{}
	// 	if err := user.QueryByKID(a.Sqlx, suser.GetUserID()); err != nil {
	// 		return err
	// 	}
	// 	account := &schema.SigninGpaAccount{}
	// 	if err := account.DeleteByUserAndKind(a.Sqlx, user.ID, int(schema.ATOpenid), b.KID); err != nil {
	// 		logger.Errorf(c, logger.ErrorWW(err)) // 无法解绑用户
	// 		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-UNBIND-NOACC", Other: "无法获取当前账户"})
	// 	}
	// 	return nil
	// }
	return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-UNBIND-NONE", Other: "无效第三方登陆"})
}
