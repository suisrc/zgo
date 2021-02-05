package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
)

// Connect system
type Connect struct {
	gpa.GPA
	Auther        auth.Auther
	SigninService *service.Signin
	CasbinAuther  *module.CasbinAuther
	SigninAPI     *Signin
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Connect) Register(r gin.IRouter) {
	c := r.Group("connect")

	c.GET("oauth2/authorize", a.authorize) // OAUTH2认证时候使用

	c.GET("oauth2/signin", a.signin)          // OAUTH2登陆使用GET请求
	c.GET("oauth2/signout", a.signout)        // OAUTH2登陆使用GET请求
	c.GET("oauth2/signin/:p", a.signin)       // OAUTH2登陆使用GET请求
	c.GET("oauth2/signin/:p/:g", a.signin)    // OAUTH2登陆使用GET请求
	c.GET("oauth2/signin/:p/:g/:w", a.signin) // OAUTH2登陆使用GET请求

	c.POST("3rd/token", a.SigninAPI.signin)         // 获取新的访问令牌
	c.GET("3rd/token/get", a.SigninAPI.token3get)   // 获取新的访问令牌
	c.GET("3rd/token/refresh", a.SigninAPI.refresh) // 获取新的访问令牌
}

// 需要完成oauth2认证
func (a *Connect) authorize(c *gin.Context) {
	helper.ResSuccess(c, "OAUTH2认证功能未开放")
}

// oauth2登录
func (a *Connect) signin(c *gin.Context) {
	// 解析参数
	body := schema.SigninOfOAuth2{}
	if err := c.ShouldBindUri(&body); err != nil {
		// c.Param()
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 处理填充参数， 再使用param参数时候， 重置空
	helper.IfExec(body.Platform == "0", func() { body.Platform = "" }, nil)
	helper.IfExec(body.OrgCode == "0", func() { body.OrgCode = "" }, nil)
	helper.IfExec(body.WebToken == "0", func() { body.WebToken = "" }, nil)
	// Form Or Query
	if err := helper.ParseForm(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 执行登录， 验证用户
	user, err := a.SigninService.OAuth2(c, &body, a.SigninAPI.LastSignIn)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 执行登录， 生成令牌
	a.SigninAPI.signin2(c, user)
}

// oauth2登出
func (a *Connect) signout(c *gin.Context) {
	a.SigninAPI.signout(c)
}
