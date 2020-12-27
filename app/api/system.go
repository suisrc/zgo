package api

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/helper"
)

// System system
type System struct {
	gpa.GPA
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *System) Register(r gin.IRouter) {
	r.GET("pub/system/info", a.getSystemInfo)
}

// signin godoc
// @Tags sign
// @Summary System
// @Description 系统信息
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /pub/system/info [get]
func (a *System) getSystemInfo(c *gin.Context) {
	helper.ResSuccess(c, helper.H{
		"status":    "ok",
		"key":       "kratos",
		"name":      "KA单点登录系统",
		"copyright": "Copyright © 2020 mecoolchina.com",
		"beian":     "沪ICP备16008893号",
		"iconfonts": []string{"//at.alicdn.com/t/font_1866669_1ybdy8kelkq.js"},
		"links": []helper.H{
			{
				"key":   "plus",
				"title": "普乐师（上海）数字科技股份有限公司",
				"href":  "http://www.mecoolchina.com",
			},
		},
	})
}
