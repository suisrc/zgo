package api

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/helper"
)

// Auth auth
type Auth struct {
	CasbinAuther *service.CasbinAuther
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Auth) Register(r gin.IRouter) {
	uaz := a.CasbinAuther.UserAuthCasbinMiddlewareByOrigin(fixRequestHeaderParam)
	uax := a.CasbinAuther.UserAuthBasicMiddleware()
	r.GET("authz", uaz, a.authorize)
	r.GET("authx", uax, a.authorize)
	r.GET("authz/clear", uax, a.clear)
}

// fixRequestHeaderParam 修复请求头的内容
func fixRequestHeaderParam(c *gin.Context, k string) (string, error) {
	// X-Reqeust-Origin-Path
	// nginx.ingress.kubernetes.io/configuration-snippet: |
	//   proxy_set_header X-Request-Origin-Host $host;
	//   proxy_set_header X-Request-Origin-Path $request_uri;
	//   proxy_set_header X-Request-Origin-Method $request_method;
	value := c.GetHeader(k)
	if value == "" {
		if k == helper.XReqOriginHostKey {
			return "default", nil
		}
		return "", errors.New("invalid " + k)
	} else if offset := strings.IndexRune(value, '?'); offset > 0 {
		value = value[:offset]
	}
	return value, nil
}

func (a *Auth) clear(c *gin.Context) {
	user, _ := helper.GetUserInfo(c)

	if user.GetOrgCode() == schema.PlatformCode && user.GetOrgAdmin() == schema.SuperUser {
		org := c.Request.FormValue("org")
		a.CasbinAuther.ClearEnforcer(org == "", org)
		helper.ResSuccess(c, "ok")
	} else {
		helper.ResSuccess(c, "error")
	}

}

// @Param Authorization header string true "Bearer token"

// Authorize godoc
// @Tags auth
// @Summary Authorize
// @Description 授权接口
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /authz [get]
func (a *Auth) authorize(c *gin.Context) {
	// 权限判断有UserAuthCasbinMiddleware完成
	// 仅仅返回正常结果即可

	// 如果通过验证， 当前用户是一定登录的
	user, _ := helper.GetUserInfo(c)

	h := c.Writer.Header()
	h.Set("X-Request-Z-Token-Kid", user.GetTokenID())

	if user.GetAccount() != "" {
		if acc, usr, err := service.DecryptAccountWithUser(c, user.GetAccount(), user.GetTokenID()); err == nil {
			h.Set("X-Request-Z-Account", strconv.Itoa(acc))
			h.Set("X-Request-Z-User-Idx", strconv.Itoa(usr))
		}
	}

	h.Set("X-Request-Z-User-Kid", user.GetUserID())
	h.Set("X-Request-Z-User-Name", url.QueryEscape(user.GetUserName()))
	h.Set("X-Request-Z-User-Roles", strings.Join(user.GetUserRoles(), ";"))

	h.Set("X-Request-Z-Org-Code", user.GetOrgCode())
	h.Set("X-Request-Z-Org-Admin", user.GetOrgAdmin())
	h.Set("X-Request-Z-Org-Usrid", user.GetOrgUsrID())
	h.Set("X-Request-Z-Org-Appid", user.GetOrgAppID())

	h.Set("X-Request-Z-Domain", user.GetDomain())
	h.Set("X-Request-Z-Issuer", user.GetIssuer())
	h.Set("X-Request-Z-Audience", user.GetAudience())

	h.Set("X-Request-Z-Xip", helper.GetHostIP(c))
	helper.ResSuccess(c, "ok")
}
