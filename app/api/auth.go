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

	// XReqUserKey         = "X-Request-User-Kid"     // user kid
	// XReqRoleKey         = "X-Request-Role-Kid"     // role kid
	// XReqDomainKey       = "X-Request-Domain"       // domain
	// XReqOrganizationKey = "X-Request-Organization" // Organization
	// XReqAccountKey      = "X-Request-Account"      // account
	// XReqUserIdxKey      = "X-Request-User-Xid"     // user index id
	// XreqUser3rdKey      = "X-Request-User-Tid"     // user third id (application)
	// XReqRoleOrgKey      = "X-Request-Role-Org"     // role organization kid
	// XReqZgoKey          = "X-Request-Zgo-Uri"      // 由于前置授权无需应用间绑定， 如果需要执行必要通信，可以获取通信地址

	h.Set(helper.XReqUserKey, user.GetUserID())
	h.Set(helper.XReqRoleKey, user.GetRoleID())
	h.Set(helper.XReqAccountKey, user.GetAccountID())
	h.Set(helper.XreqUserNamKey, user.GetUserName())
	h.Set(helper.XReqUserIdxKey, user.GetXID())
	h.Set(helper.XreqUser3rdKey, user.GetTID())
	h.Set(helper.XReqDomainKey, "nil")       // 平台
	h.Set(helper.XReqOrganizationKey, "nil") // 平台 LCOAL-PM-00
	h.Set(helper.XReqRoleOrgKey, "nil")      // 平台

	h.Set(helper.XReqZgoKey, helper.GetHostIP(c))

	helper.ResSuccess(c, "ok")
}
