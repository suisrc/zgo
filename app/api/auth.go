package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
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
	uac := middleware.UserAuthCasbinMiddleware(a.Auther, a.Enforcer)

	r.GET("authz", uac, a.authorize)
	// r.GET("authz/signin", uac, a.signin)
	// r.GET(middleware.JoinPath(config.C.HTTP.ContextPath, "authz"), uac, a.authorize)
}

// Register 主路由必须包含UAC内容
func (a *Auth) Register(r gin.IRouter) {
	r.GET("authz", a.authorize)
	r.GET("authz/signin", a.signin)
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
	helper.ResSuccess(c, "ok")
}

// Signin godoc
// @Tags auth
// @Summary Authorize
// @Description 授权接口
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /authz/signin [get]
func (a *Auth) signin(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}
