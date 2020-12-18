package manager

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// Role 角色管理器
type Role struct {
	gpa.GPA
}

// Register 主路由必须包含UAC内容
func (a *Role) Register(r gin.IRouter) {
	u := r.Group("roles")

	u.GET("list", a.RolesList)
	u.GET("size", a.RolesSize)
}

// RolesList 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/roles/list [get]
func (a *Role) RolesList(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

// RolesSize 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/roles/size [get]
func (a *Role) RolesSize(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
