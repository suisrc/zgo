package manager

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// Menu 菜单管理器
type Menu struct {
	gpa.GPA
}

// Register 主路由必须包含UAC内容
func (a *Menu) Register(r gin.IRouter) {
	u := r.Group("menu")

	u.GET("list", a.MenusList)
	u.GET("size", a.MenusSize)
}

// MenusList 获取用户列表
func (a *Menu) MenusList(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

// MenusSize 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/account/list [get]
func (a *Menu) MenusSize(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
