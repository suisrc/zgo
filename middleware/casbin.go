package middleware

import (
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	conf := config.C.Casbin
	if !conf.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		u, ok := helper.GetUserInfo(c)
		if !ok {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		}

		//i := u.GetUserID()
		r := u.GetRoleID()
		a := u.GetAudience()
		d := c.Request.URL.Host
		p := c.Request.URL.Path
		i := helper.GetClientIP(c)
		m := c.Request.Method
		if b, err := enforcer.Enforce(r, a, d, p, i, m); err != nil {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		} else if !b {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		}
		c.Next()
	}
}
