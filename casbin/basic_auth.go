package casbin

import (
	"github.com/suisrc/zgo/auth"
	"github.com/suisrc/zgo/helper"
	"github.com/suisrc/zgo/middleware"

	"github.com/gin-gonic/gin"
)

// UseAuthBasicMiddleware 用户授权中间件, 只判定登录权限
func (a *Auther) UseAuthBasicMiddleware(skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if middleware.SkipHandler(c, skippers...) {
			c.Next() // 需要跳过权限验证的uri内容
			return
		}

		user, err := a.Implor.GetAuther().GetUserInfo(c, "")
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
