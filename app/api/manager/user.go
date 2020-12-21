package manager

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// User 用户管理器
type User struct {
	gpa.GPA
}

// Register 主路由必须包含UAC内容
func (a *User) Register(r gin.IRouter) {
	u := r.Group("users")

	u.GET("list", a.UsersList)
	u.GET("size", a.UsersSize)
}

// UsersList 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/users/list [get]
func (a *User) UsersList(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

// UsersSize 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/users/size [get]
func (a *User) UsersSize(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
