package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/helper"
)

// User 用户管理器
type User struct {
	service.GPA
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
