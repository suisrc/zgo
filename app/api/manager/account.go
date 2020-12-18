package manager

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// Account 账户管理器
type Account struct {
	gpa.GPA
}

// Register 主路由必须包含UAC内容
func (a *Account) Register(r gin.IRouter) {
	u := r.Group("account")

	u.GET("list", a.AccountsList)
	u.GET("size", a.AccountsSize)
}

// AccountsList 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/account/list [get]
func (a *Account) AccountsList(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

// AccountsSize 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/account/size [get]
func (a *Account) AccountsSize(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
