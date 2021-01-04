package api

import (
	"errors"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/middleware"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
)

// Auth auth
type Auth struct {
	Enforcer *casbin.SyncedEnforcer
	Auther   auth.Auther
}

// RegisterWithUAC 注册路由,认证接口特殊,需要独立注册
func (a *Auth) RegisterWithUAC(r gin.IRouter) {
	uac := middleware.UserAuthCasbinMiddlewareByPathFunc(a.Auther, a.Enforcer, func(c *gin.Context, k string) (string, error) {
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
	})

	r.GET("authz", uac, a.authorize)
	// r.GET("authz/signin", uac, a.signin)
	// r.GET(middleware.JoinPath(config.C.HTTP.ContextPath, "authz"), uac, a.authorize)
}

// Register 主路由必须包含UAC内容
func (a *Auth) Register(r gin.IRouter) {
	r.GET("authz", a.authorize)
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
	user, exist := helper.GetUserInfo(c)
	if !exist {
		// 未登录
		helper.ResError(c, &helper.ErrorModel{
			Status:   200,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-AUTHORIZE-USERNOEXIST",
				Other: "登录用户不存在",
			},
		})
		return
	}

	h := c.Writer.Header()
	h.Set("X-Request-Z-Token-Kid", user.GetTokenID())
	h.Set("X-Request-Z-User-Kid", user.GetUserID())
	h.Set("X-Request-Z-User-Nam"， user.GetUserName())
	h.Set("X-Request-Z-Role-Kid"， user.GetUserRole())
	h.Set("X-Request-Z-User-Xid"， user.GetXidxID())
	h.Set("X-Request-Z-Account"， user.GetAccountID())
	h.Set("X-Request-Z-User-Tid", user.GetT3rdID())
	h.Set("X-Request-Z-Client-Kid"， user.GetClientID())
	h.Set("X-Request-Z-Domain", user.GetDomain())
	h.Set("X-Request-Z-Issuer"， user.GetIssuer())
	h.Set("X-Request-Z-Audience"， user.GetAudience())
	h.Set("X-Request-Z-Org-Code"， user.GetOrgCode())
	h.Set("X-Request-Z-Org-Domain", user.GetOrgDomain())
	h.Set("X-Request-Z-Org-Admin", user.GetOrgAdmin())

	//h.Set("X-Request-Z-Org-Code"， "ORGCM3558")

	h.Set(helper.XReqZgoKey, helper.GetHostIP(c))
	helper.ResSuccess(c, "ok")
}
