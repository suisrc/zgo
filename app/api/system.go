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
	r.GET("/pub/3rd/apps", a.getAppList)
	//r.GET("signin/mfa", a.signinMFA)
	//r.POST("signup", a.signup) // 注册
	//r.GET("signin/oauth2/:kid", a.oauth2) // OAUTH2登陆使用GET请求
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

// signin godoc
// @Tags sign
// @Summary Use3rd
// @Description 系统信息
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /pub/3rd/apps [get]
func (a *System) getAppList(c *gin.Context) {
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

//==================================================================================================================
//==================================================================================================================
//==================================================================================================================

// oauth2 godoc
// @Tags sign
// @Summary OAuth2
// @Description 第三方授权登陆
// @Accept  json
// @Produce  json
// @Param kid path string true "平台KID"
// @Param result query string false "返回值类型, 比如: json"
// @Param redirect query string false "redirect"
// @Success 200 {object} helper.Success
// @Router /signin/oauth2/{kid} [get]
//func (a *Signin) oauth2(c *gin.Context) {
// 解析参数
//body := schema.SigninOfOAuth2{}
//if err := helper.ParseQuery(c, &body); err != nil {
//	helper.FixResponse406Error(c, err, func() {
//		logger.Errorf(c, logger.ErrorWW(err))
//	})
//	return
//}
//
//// 执行登录
//user, err := a.SigninService.OAuth2(c, &body, a.last)
//if err != nil {
//	helper.FixResponse500Error(c, err, func() {
//		logger.Errorf(c, logger.ErrorWW(err))
//	})
//	return
//}
//token, usr, err := a.Auther.GenerateToken(c, user)
//if err != nil {
//	helper.FixResponse500Error(c, err, func() {
//		logger.Errorf(c, logger.ErrorWW(err))
//	})
//	return
//}
//
//// 登陆日志
//a.log(c, usr, token, "oauth2", token.GetRefreshToken())
//if body.Redirect != "" {
//	// 需要重定向跳转
//	redirect, err := url.QueryUnescape(body.Redirect)
//	if err != nil {
//		helper.FixResponse500Error(c, err, func() {
//			logger.Errorf(c, logger.ErrorWW(err))
//		})
//		return
//	}
//	if strings.IndexRune(redirect, '?') <= 0 {
//		redirect += "?"
//	}
//	if endc := redirect[len(redirect)-1:]; endc != "?" && endc != "&" {
//		redirect += "&"
//	}
//	redirect += "access_token=" + token.GetAccessToken()
//	redirect += "&expires_at=" + strconv.Itoa(int(token.GetExpiresAt()))
//	redirect += "&expires_in=" + strconv.Itoa(int(token.GetExpiresAt()-time.Now().Unix()))
//	redirect += "&refresh_token=" + token.GetRefreshToken()
//	redirect += "&refresh_expires=" + strconv.Itoa(int(token.GetRefreshExpAt()))
//	redirect += "&token_type=Bearer"
//	redirect += "&trace_id=" + helper.GetTraceID(c)
//	// 重定向到登陆页面
//	c.Redirect(303, redirect)
//	return
//}
//
//// 登陆结果
//result := schema.SigninResult{
//	TokenStatus:  "ok",
//	TokenType:    "Bearer",
//	TokenID:      token.GetTokenID(),
//	AccessToken:  token.GetAccessToken(),
//	ExpiresAt:    token.GetExpiresAt(),
//	ExpiresIn:    token.GetExpiresAt() - time.Now().Unix(),
//	RefreshToken: token.GetRefreshToken(),
//	RefreshExpAt:   token.GetRefreshExpAt(),
//}
//
//// 记录登陆
//// 返回正常结果即可
//helper.ResSuccess(c, &result)
//}

// Signup godoc
// @Tags sign
// @Summary Signup
// @Description 登陆
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /signup [post]
//func (a *Signin) signup(c *gin.Context) {
//	helper.ResSuccess(c, "功能未开放")
//}
