package api

import (
	"time"

	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
)

// User 用户管理器
type User struct {
	gpa.GPA
	UserService service.User
	Auther      auth.Auther // 令牌控制
}

// Register 注册接口
func (a *User) Register(r gin.IRouter) {
	user := r.Group("user")

	current := user.Group("current")
	{
		current.GET("", a.current)
		current.GET("access", a.access)
		current.GET("notices", a.notices)
	}
	oauth2 := user.Group("oauth2")
	{
		oauth2.GET("bind", a.bindOAuth2)
		oauth2.GET("unbind", a.unbindOAuth2)
	}
	passwd := user.Group("passwd")
	{
		passwd.POST("change", a.changePasswd)
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
func (a *User) current(c *gin.Context) {
	// 验证登录信息
	user, err := a.Auther.GetUserInfo(c)
	// 如果通过验证， 当前用户是一定登录的
	// user, exist := helper.GetUserInfo(c)
	if err != nil || user == nil /*!exist*/ {
		// 未登录
		helper.ResError(c, &helper.ErrorModel{
			Status:   200,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-AUTHORIZE-USERNOEXIST",
				Other: "登录用户不存在",
			},
		})
		return
	}

	// userid?: string; // 用户ID
	// avatar?: string; // 头像
	// name?: string; // 名称
	// system?: string; // 该字段主要是有前端给出,用以记录使用, 不同system带来的access也是不同的
	// createAt?: number; // 获取当前信息的时间
	helper.ResSuccess(c, helper.H{
		"userid":   user.GetUserID(),
		"name":     user.GetUserName(),
		"system":   "LHDG2",
		"createAt": time.Now(),
	})
}

/**
 * 动态验证用户权限问题
 */
// @Tags user
// @Summary 动态验证用户权限问题
// @Description 动态验证用户权限问题
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/current/access [get]
func (a *User) access(c *gin.Context) {
	// helper.ResSuccess(c, "ok")

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowNone,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}

/**
 * 查询当前用户消息
 */
// @Tags user
// @Summary 查询当前用户消息
// @Description 查询当前用户消息
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /user/current/notices [get]
func (a *User) notices(c *gin.Context) {
	// helper.ResSuccess(c, "ok")

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowNone,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
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
func (a *User) bindOAuth2(c *gin.Context) {
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
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param kid query string true "平台KID"
// @Success 200 {object} helper.Success
// @Router /user/oauth2/unbind [get]
func (a *User) unbindOAuth2(c *gin.Context) {
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

// @Tags user
// @Summary passwd
// @Description 修改密码
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param kid query string true "平台KID"
// @Success 200 {object} helper.Success
// @Router /user/passwd/change [post]
func (a *User) changePasswd(c *gin.Context) {
	// helper.ResSuccess(c, "ok")

	helper.ResError(c, &helper.ErrorModel{
		Status:   200,
		ShowType: helper.ShowWarn,
		ErrorMessage: &i18n.Message{
			ID:    "ERR-INTERFACE-NOTOPEN",
			Other: "功能接口为开放",
		},
	})
}
