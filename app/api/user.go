package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
)

// User 用户管理器
type User struct {
	gpa.GPA
	UserService service.User
}

// Register 注册接口
func (a *User) Register(r gin.IRouter) {
	user := r.Group("user")

	current := user.Group("current")
	{
		current.GET("", a.userCurrent)
		current.GET("access", a.userCurrentAccess)
		current.GET("notices", a.userCurrentNotices)
	}
	oauth2 := user.Group("oauth2")
	{
		oauth2.GET("bind", a.bindOAuth2Account)
		oauth2.GET("unbind", a.unbindOAuth2Account)
	}
}

/**
 * 查询当前用户信息
 * 一般在用户收起登陆，或者首次打开页面时候触发
 * 只有基本信息
 */
// @Tags user
// @Summary 查询当前用户信息
// @Description 查询当前用户信息
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/current [get]
func (a *User) userCurrent(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

/**
 * 动态验证用户权限问题
 */
// @Tags user
// @Summary 查询当前用户信息
// @Description 查询当前用户信息
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/current/access [get]
func (a *User) userCurrentAccess(c *gin.Context) {

	helper.ResSuccess(c, "ok")
}

/**
 * 查询当前用户信息
 */
// @Tags user
// @Summary 查询当前用户信息
// @Description 查询当前用户信息
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/current/notices [get]
func (a *User) userCurrentNotices(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

// @Tags user
// @Summary Bind
// @Description 绑定第三方账户
// @Accept  json
// @Produce  json
// @Param kid query string true "平台KID"
// @Param result query string false "返回值类型, 比如: json"
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/oauth2/bind [get]
func (a *User) bindOAuth2Account(c *gin.Context) {
	// 解析参数
	body := schema.SigninOfOAuth2{}
	if err := helper.ParseQuery(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 执行绑定
	err := a.UserService.Bind(c, &body)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	helper.ResSuccess(c, "ok")
}

// @Tags user
// @Summary Unbind
// @Description 解绑第三方账户
// @Accept  json
// @Produce  json
// @Param kid query string true "平台KID"
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/oauth2/unbind [get]
func (a *User) unbindOAuth2Account(c *gin.Context) {
	// 解析参数
	body := schema.SigninOfOAuth2{}
	if err := helper.ParseQuery(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 执行绑定
	err := a.UserService.Unbind(c, &body)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	helper.ResSuccess(c, "ok")
}
