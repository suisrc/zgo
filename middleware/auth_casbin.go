package middleware

import (
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

var (
	// NoSignin 未登陆
	NoSignin = "nosignin"
	// NoRole 无角色
	NoRole = "norole"
)

// UserAuthCasbinMiddleware 用户授权中间件
func UserAuthCasbinMiddleware(auther auth.Auther, enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return EmptyMiddleware()
	}
	conf := config.C.Casbin

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		user, err := auther.GetUserInfo(c)
		if !conf.Enable {
			// casbin禁用, 只判定是否登陆
			if err != nil {
				// 获取登陆信息异常
				if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
					helper.ResError(c, helper.Err401Unauthorized)
					return
				}
				helper.ResError(c, helper.Err400BadRequest)
				return
			}
			helper.SetUserInfo(c, user)
			c.Next()
			return
		}

		// 需要执行casbin授权
		var r, a string               // 角色, jwt授权方
		erm := helper.Err403Forbidden // casbin验证失败后返回的异常

		if err != nil {
			if err == auth.ErrNoneToken && conf.NoSignin {
				r = NoSignin                    // 用户未登陆,且允许执行未登陆认证
				erm = helper.Err401Unauthorized // 替换403异常,因为当前用户未登陆
			} else if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
				helper.ResError(c, helper.Err401Unauthorized) // 无有效登陆用户
				return
			} else {
				helper.ResError(c, helper.Err400BadRequest) // 解析jwt令牌出现未知错误
				return
			}
		} else {
			r = user.GetRoleID() // 请求角色
			if r == "" && conf.NoRole {
				r = NoRole // 用户无角色,且允许执行无角色认证
			} else {
				helper.ResError(c, helper.Err403Forbidden) // 无角色,禁止访问
				return
			}
			a = user.GetAudience() // jwt授权方
		}

		d := c.Request.URL.Host    // 请求域名
		p := c.Request.URL.Path    // 请求路径
		i := helper.GetClientIP(c) // 客户端IP
		m := c.Request.Method      // 请求方法
		if b, err := enforcer.Enforce(r, a, d, p, i, m); err != nil {
			logger.Errorf(c, err.Error()) // 授权发生异常
			helper.ResError(c, erm)
			return
		} else if !b {
			helper.ResError(c, erm)
			return
		}

		if user != nil {
			helper.SetUserInfo(c, user)
		}
		c.Next()
	}
}
