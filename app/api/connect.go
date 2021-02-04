package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
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
	c.GET("oauth2/signin", a.signin)       // OAUTH2登陆使用GET请求
	// c.GET("oauth2/signin/:p", a.signin)       // OAUTH2登陆使用GET请求
	// c.GET("oauth2/signin/:p/:g", a.signin)    // OAUTH2登陆使用GET请求
	// c.GET("oauth2/signin/:p/:g/:w", a.signin) // OAUTH2登陆使用GET请求

	c.POST("3rd/token", a.Signin.signin)         // 获取新的访问令牌
	c.GET("3rd/token/get", a.Signin.token3get)   // 获取新的访问令牌
	c.GET("3rd/token/refresh", a.Signin.refresh) // 获取新的访问令牌
}

// 需要完成oauth2认证
func (a *Connect) authorize(c *gin.Context) {
	helper.ResSuccess(c, "功能未开放")
}

// oauth2登录
func (a *Connect) signin(c *gin.Context) {
	// 解析参数
	b := schema.SigninOfOAuth2{}
	if err := helper.ParseJSON(c, &b); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	a.fixSigninOfOAuth2(c, &b)
}

// fix SigninOfOAuth2
func (a *Connect) fixSigninOfOAuth2(c *gin.Context, b *schema.SigninOfOAuth2) {
	if b.Platform == "x" {
		b.Platform = ""
	}
	if b.OrgCode == "x" {
		b.OrgCode = ""
	}
	if b.WebToken == "x" {
		b.WebToken = ""
	}
}
