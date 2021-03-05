package middleware

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/suisrc/zgo/auth"
	"github.com/suisrc/zgo/helper"

	"github.com/gin-gonic/gin"
)

// UseAuthServerMiddleware 用户授权中间件
func UseAuthServerMiddleware(c *gin.Context, gf func(*gin.Context, auth.UserInfo) (aid, uid int64, err error)) {
	// 如果通过验证， 当前用户是一定登录的
	user, _ := helper.GetUserInfo(c)

	h := c.Writer.Header()
	h.Set("X-Request-Z-Token-Kid", user.GetTokenID())
	h.Set("X-Request-Z-Token-Pid", user.GetTokenPID())

	if user.GetAccount() != "" && gf != nil {
		if acc, usr, err := gf(c, user); err == nil {
			h.Set("X-Request-Z-Account-Id", strconv.Itoa(int(acc)))
			if usr > 0 {
				h.Set("X-Request-Z-User-Id", strconv.Itoa(int(usr)))
			}
		}
	}

	h.Set("X-Request-Z-Account", user.GetAccount())   // 账户名
	h.Set("X-Request-Z-Account1", user.GetAccount1()) // 账户名1
	h.Set("X-Request-Z-Account2", user.GetAccount2()) // 账户名2
	h.Set("X-Request-Z-User-Kid", user.GetUserID())   // 用户唯一标识
	h.Set("X-Request-Z-User-Name", url.QueryEscape(user.GetUserName()))
	h.Set("X-Request-Z-User-Roles", strings.Join(user.GetUserRoles(), ";"))
	h.Set("X-Request-Z-Org-Code", user.GetOrgCode())
	h.Set("X-Request-Z-Org-Admin", user.GetOrgAdmin())
	h.Set("X-Request-Z-Org-Usrid", user.GetOrgUsrID())
	h.Set("X-Request-Z-Agent", user.GetAgent())
	h.Set("X-Request-Z-Scope", user.GetScope())
	h.Set("X-Request-Z-Domain", user.GetDomain())
	h.Set("X-Request-Z-Issuer", user.GetIssuer())
	h.Set("X-Request-Z-Audience", user.GetAudience())
	h.Set("X-Request-Z-Xip", helper.GetHostIP(c))

	c.Next()
}
