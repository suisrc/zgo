package api

import (
	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
)

// Signin signin
type Signin struct {
	Auther        auth.Auther
	SigninService service.Signin
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Signin) Register(r gin.IRouter) {
	// sign 开头的路由会被全局casbin放行
	r.POST("signin", a.signin)         // 登陆必须是POST请求
	r.POST("signin/{:kid}", a.signin2) // 登陆必须是POST请求

	// ua := middleware.UserAuthMiddleware(a.Auther)
	// r.GET("signout", ua, a.signout)
	r.GET("signout", a.signout)
	r.GET("signin/refresh", a.refresh)

	r.POST("signup", a.signup)         // 注册
	r.POST("signup/{:kid}", a.signup2) // 注册
}

// Signin godoc
// @Tags sign
// @Summary Signin
// @Description 登陆
// @Accept  json
// @Produce  json
// @Param item body schema.SigninBody true "SigninBody Info"
// @Success 200 {object} helper.Success
// @Router /signin [post]
func (a *Signin) signin(c *gin.Context) {
	// 解析参数
	body := schema.SigninBody{}
	if err := helper.ParseJSON(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, err.Error())
		})
		return
	}
	// 执行登录
	user, err := a.SigninService.Signin(c, &body)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, err.Error())
		})
		return
	}
	token, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, err.Error())
		})
		return
	}

	result := schema.SigninResult{
		Status:  "ok",
		Token:   token.GetAccessToken(),
		Expired: token.GetExpiresAt(),
		//Expired: token.GetExpiresAt() - time.Now().Unix(),
	}
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}
func (a *Signin) signin2(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

// Signout godoc
// @Tags sign
// @Summary Signin
// @Description 登出
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /signout [get]
func (a *Signin) signout(c *gin.Context) {
	// 返回正常结果即可
	// user, b := helper.GetUserInfo(c)

	// 确定登陆用户的身份
	user, err := a.Auther.GetUserInfo(c)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		}
		helper.ResError(c, helper.Err400BadRequest)
		return
	}

	// 执行登出
	if err := a.Auther.DestroyToken(c, user); err != nil {
		helper.ResError(c, helper.Err400BadRequest)
		return
	}

	helper.ResSuccess(c, "ok")
}

// Refresh godoc
// @Tags sign
// @Summary Refresh
// @Description 刷新令牌
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /signin/refresh [get]
func (a *Signin) refresh(c *gin.Context) {
	// 确定登陆用户的身份
	user, err := a.Auther.GetUserInfo(c)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		}
		helper.ResError(c, helper.Err400BadRequest)
		return
	}
	token, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, err.Error())
		})
		return
	}

	result := schema.SigninResult{
		Status:  "ok",
		Token:   token.GetAccessToken(),
		Expired: token.GetExpiresAt(),
		//Expired: token.GetExpiresAt() - time.Now().Unix(),
	}
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}

// Signup godoc
// @Tags sign
// @Summary Signup
// @Description 登陆
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /signup [post]
func (a *Signin) signup(c *gin.Context) {
	helper.ResSuccess(c, "功能为开放")
}

// 注册
func (a *Signin) signup2(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

// 绑定
func (a *Signin) signbind(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}
