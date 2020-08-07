package middleware

import (
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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
		if err != nil {
			if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
				helper.ResError(c, &helper.Err401Unauthorized)
				return
			}
			helper.ResError(c, &helper.Err400BadRequest)
			return
		}

		if conf.Enable {
			r := user.GetRoleID()
			a := user.GetAudience()
			d := c.Request.URL.Host
			p := c.Request.URL.Path
			i := helper.GetClientIP(c)
			m := c.Request.Method
			if b, err := enforcer.Enforce(r, a, d, p, i, m); err != nil {
				logger.Errorf(c, err.Error())
				helper.ResError(c, &helper.Err403Forbidden)
				return
			} else if !b {
				helper.ResError(c, &helper.Err403Forbidden)
				return
			}
		}

		helper.SetUserInfo(c, user)
		c.Next()
	}
}
