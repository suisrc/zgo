package api

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// Use3rd use3rd
type Use3rd struct {
	gpa.GPA
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Use3rd) Register(r gin.IRouter) {
	r.GET("/pub/3rd/apps", a.list)
	r.GET("/pub/3rd/signin", a.signin)
}

// signin godoc
// @Tags sign
// @Summary Use3rd
// @Description 系统信息
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /pub/3rd/apps [get]
func (a *Use3rd) list(c *gin.Context) {
	helper.ResSuccess(c, []helper.H{})
	// helper.ResSuccess(c, []helper.H{
	// 	{
	// 		"platform":  "wechat",
	// 		"appid":     "10001",
	// 		"name":      "微信10",
	// 		"title":     "微信11",
	// 		"signature": "10001",
	// 		"icon":      "iconwechat-fill",
	// 	},
	// 	{
	// 		"platform":  "github",
	// 		"appid":     "30001",
	// 		"name":      "GitHub10",
	// 		"title":     "GitHub11",
	// 		"signature": "30001",
	// 		"icon":      "icongithub-fill",
	// 	},
	// })
}

// signin godoc
// @Tags sign
// @Summary Sign
// @Description 登录
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /pub/3rd/signin [get]
func (a *Use3rd) signin(c *gin.Context) {
	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
