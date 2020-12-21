package manager

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// Gateway 网关管理器
type Gateway struct {
	gpa.GPA
}

// Register 主路由必须包含UAC内容
func (a *Gateway) Register(r gin.IRouter) {
	u := r.Group("gateway")

	u.GET("list", a.GatewaysList)
	u.GET("size", a.GatewaysSize)
}

// GatewaysList 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/gateway/list [get]
func (a *Gateway) GatewaysList(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

// GatewaysSize 获取用户列表
// @Tags manager
// @Summary 获取用户列表
// @Description 获取用户列表
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /manager/gateway/size [get]
func (a *Gateway) GatewaysSize(c *gin.Context) {

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
