package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/modules/helper"
)

// Connect system
type Connect struct {
	gpa.GPA
	Signin       *Signin
	CasbinAuther *module.CasbinAuther
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Connect) Register(r gin.IRouter) {
	c := r.Group("connect")

	c.GET("oauth2/authorize", a.authorize) // OAUTH2登陆使用GET请求

	c.POST("3rd/token", a.Signin.signin)         // 获取新的访问令牌
	c.GET("3rd/token/get", a.Signin.token3get)   // 获取新的访问令牌
	c.GET("3rd/token/refresh", a.Signin.refresh) // 获取新的访问令牌
}

// 需要完成oauth2认证
func (a *Connect) authorize(c *gin.Context) {
	helper.ResSuccess(c, "功能未开放")
}
