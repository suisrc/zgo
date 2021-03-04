package middleware

import (
	"github.com/suisrc/zgo/auth"
	"github.com/suisrc/zgo/config"
	"github.com/suisrc/zgo/helper"

	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware 用户授权中间件,废弃,请使用UserAuthCasbinMiddleware
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		user, err := a.GetUserInfo(c, "")
		if err != nil {
			if err == auth.ErrNoneToken || err == auth.ErrInvalidToken {
				helper.ResError(c, helper.Err401Unauthorized)
				return // 无有效登陆用户
			} else if err == auth.ErrExpiredToken {
				helper.ResError(c, helper.Err456TokenExpired)
				return // 访问令牌已经过期
			}
			helper.ResError(c, helper.Err500InternalServer)
			return // 解析jwt令牌出现未知错误
		}
		helper.SetUserInfo(c, user)
		c.Next()
		return
	}
}
