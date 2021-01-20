package api

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
)

// Signin signin
type Signin struct {
	gpa.GPA
	Auther        auth.Auther
	SigninService service.Signin
}

// Register 注册路由,认证接口特殊,需要独立注册
// sign 开头的路由会被全局casbin放行
func (a *Signin) Register(r gin.IRouter) {

	r.POST("signin", a.signin)            // 登录系统， 获取令牌 POST请求
	r.GET("signout", a.signout)           // 登出系统， 注销令牌（访问令牌和刷新令牌）
	r.GET("signin/refresh", a.refresh)    // 刷新令牌
	r.GET("signin/captcha", a.captcha)    // 发送验证码
	r.GET("signin/token/new", a.tokenNew) // 构建新的访问令牌
	r.GET("signin/token/get", a.tokenGet) // 获取新的访问令牌
	//r.GET("signin/mfa", a.signinMFA)
	//r.POST("signup", a.signup) // 注册
	//r.GET("signin/oauth2/:kid", a.oauth2) // OAUTH2登陆使用GET请求

}

// signin godoc
// @Tags sign
// @Summary Signin
// @Description 登陆
// @Accept  json
// @Produce  json
// @Param item body schema.SigninBody true "SigninBody Info"
// @Success 200 {object} helper.Success
// @Router /signin [post]
func (a *Signin) signin(c *gin.Context) {
	// 解析参数
	body := schema.SigninBody{}
	if err := helper.ParseJSON(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 执行登录
	user, err := a.SigninService.Signin(c, &body, a.lastSignIn)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 生成令牌
	token, usr, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 记录登录
	a.logSignIn(c, usr, token, true, "")
	// 登陆结果
	result := schema.SigninResult{
		TokenStatus:  "ok",
		TokenType:    "Bearer",
		TokenID:      token.GetTokenID(),
		AccessToken:  token.GetAccessToken(),
		ExpiresAt:    token.GetExpiresAt(),
		ExpiresIn:    token.GetExpiresAt() - time.Now().Unix(),
		RefreshToken: token.GetRefreshToken(),
		RefreshExpAt: token.GetRefreshExpAt(),
	}

	// 记录登陆
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}

// signout godoc
// @Tags sign
// @Summary Signout
// @Description 登出
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /signout [get]
func (a *Signin) signout(c *gin.Context) {
	// 返回正常结果即可
	// user, b := helper.GetUserInfo(c)

	// 确定登陆用户的身份
	user, err := a.Auther.GetUserInfo(c)
	if err != nil {
		if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
			helper.ResError(c, helper.Err401Unauthorized)
			return
		}
		helper.ResError(c, helper.Err400BadRequest)
		return
	}

	// 执行登出
	if err := a.Auther.DestroyToken(c, user); err != nil {
		helper.ResError(c, helper.Err400BadRequest)
		return
	}
	a.logSignOut(c, user, user.GetTokenID())

	helper.ResSuccess(c, "ok")
}

//==================================================================================================================

// 获取最后一次登陆信息
func (a *Signin) lastSignIn(c *gin.Context, aid int) (*schema.SigninGpaAccountToken, error) {
	if config.C.JWTAuth.LimitTime <= 0 {
		// 不使用上去签名的结果作为缓存
		return nil, nil
	}
	o2a := schema.SigninGpaAccountToken{}
	// 防止意外放生， 使用客户端IP作为影响因子
	if err := o2a.QueryByAccountAndClient(a.Sqlx, aid, helper.GetClientIP(c)); err != nil {
		if !sqlxc.IsNotFound(err) {
			// 数据库查询发生异常
			logger.Errorf(c, logger.ErrorWW(err))
			return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-DB-UNKONW", Other: "数据库发生位置异常"})
		}
	}
	if o2a.LastAt.Valid && time.Now().Unix()-o2a.LastAt.Time.Unix() < config.C.JWTAuth.LimitTime {
		// 登陆时间非常短,直接返回上次签名结果, 注意, 如果用于在很短时间从两个不同的设备登陆,会导致签发的令牌相同,并且可能会发生同时退出的问题
		// 如果需要避免上述问题,可以禁用缓存
		return nil, helper.NewSuccess(c, &schema.SigninResult{
			TokenStatus:  "ok",
			TokenType:    "Bearer",
			TokenID:      o2a.TokenID,
			AccessToken:  o2a.AccessToken.String,
			ExpiresAt:    o2a.ExpiresAt.Int64,
			ExpiresIn:    o2a.ExpiresAt.Int64 - time.Now().Unix(),
			RefreshToken: o2a.RefreshToken.String,
			RefreshExpAt: o2a.RefreshExpAt.Int64,
		})
	}
	return &o2a, nil
}

// 日志记录
func (a *Signin) logSignIn(c *gin.Context, u auth.UserInfo, t auth.TokenInfo, n bool, delay string) {
	// c.SetCookie("signin", u.GetTokenID(), -1, "", u.GetAudience(), false, false) // 标记登陆信息

	// aid, _ := strconv.Atoi(u.GetAccount())
	aid, _, err := service.DecryptAccountWithUser(c, u.GetAccount(), u.GetTokenID())
	if err != nil {
		return
	}
	// cid, cok := helper.GetCtxValueToString(c, helper.ResJwtKey)
	o2a := schema.SigninGpaAccountToken{
		TokenID:      u.GetTokenID(),
		AccountID:    aid,
		OrgCode:      sql.NullString{Valid: u.GetOrgCode() != "", String: u.GetOrgCode()},
		DelayToken:   sql.NullString{Valid: delay != "", String: delay},
		DelayExpAt:   sql.NullInt64{Valid: delay != "", Int64: time.Now().Unix() + 300},
		AccessToken:  sql.NullString{Valid: t.GetAccessToken() != "", String: t.GetAccessToken()},
		ExpiresAt:    sql.NullInt64{Valid: t.GetExpiresAt() > 0, Int64: t.GetExpiresAt()},
		RefreshToken: sql.NullString{Valid: t.GetRefreshToken() != "", String: t.GetRefreshToken()},
		RefreshExpAt: sql.NullInt64{Valid: t.GetRefreshExpAt() > 0, Int64: t.GetRefreshExpAt()},
		LastIP:       sql.NullString{Valid: true, String: helper.GetClientIP(c)},
		LastAt:       sql.NullTime{Valid: true, Time: time.Now()},
	}
	if err := o2a.UpdateAndSaveByTokenKID(a.Sqlx, !n); err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
	}
}

// 日志记录
func (a *Signin) logSignOut(c *gin.Context, u auth.UserInfo, t string) {
	// 销毁刷新令牌
	o2a := schema.SigninGpaAccountToken{
		TokenID:      u.GetTokenID(),
		RefreshExpAt: sql.NullInt64{Valid: true, Int64: 0},
	}
	if err := o2a.UpdateAndSaveByTokenKID(a.Sqlx, true); err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
	}
}

//==================================================================================================================

// refresh godoc
// @Tags sign
// @Summary Refresh
// @Description 刷新令牌
// @Accept  json
// @Produce  json
// @Param refresh_token query string true "刷新令牌"
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /signin/refresh [get]
func (a *Signin) refresh(c *gin.Context) {
	o2a := a.getSigninGpaAccountToken(c)
	if o2a == nil {
		// 无法处理， 结束
		return
	}
	if config.C.JWTAuth.LimitTime > 0 && o2a.LastAt.Valid && time.Now().Unix()-o2a.LastAt.Time.Unix() < config.C.JWTAuth.LimitTime {
		// 登陆时间非常短,直接返回上次签名结果, 注意, 如果用于在很短时间从两个不同的设备登陆,会导致签发的令牌相同,并且可能会发生同时退出的问题
		// 如果需要避免上述问题,可以禁用缓存
		result := schema.SigninResult{
			TokenStatus:  "ok",
			TokenType:    "Bearer",
			TokenID:      o2a.TokenID,
			AccessToken:  o2a.AccessToken.String,
			ExpiresAt:    o2a.ExpiresAt.Int64,
			ExpiresIn:    o2a.ExpiresAt.Int64 - time.Now().Unix(),
			RefreshToken: o2a.RefreshToken.String,
			RefreshExpAt: o2a.RefreshExpAt.Int64,
		}
		helper.ResSuccess(c, &result)
		return
	}
	// 通过刷新令牌生成新令牌
	token, user, err := a.Auther.RefreshToken(c, o2a.AccessToken.String, func(usrInfo auth.UserInfo, expIn int) error {
		if time.Now().Unix() > o2a.RefreshExpAt.Int64 {
			// 刷新令牌已经过期， 无法执行刷新
			// return errors.New("token is expired")
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-EXPIRED", Other: "刷新令牌过期"})
		}
		return nil
	})
	// 刷新令牌放生了异常， 直接结束
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	a.logSignIn(c, user, token, false, "")
	// 登录结果
	result := schema.SigninResult{
		TokenStatus:  "ok",
		TokenType:    "Bearer",
		TokenID:      token.GetTokenID(),
		AccessToken:  token.GetAccessToken(),
		ExpiresAt:    token.GetExpiresAt(),
		ExpiresIn:    token.GetExpiresAt() - time.Now().Unix(),
		RefreshToken: token.GetRefreshToken(),
		RefreshExpAt: token.GetRefreshExpAt(),
	}
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}

// 获取旧的访问令牌
func (a *Signin) getSigninGpaAccountToken(c *gin.Context) *schema.SigninGpaAccountToken {
	// 需要注意, 刷新令牌只有一次有效
	rid := c.Request.FormValue("refresh_token")
	if rid == "" {
		helper.ResError(c, helper.Err401Unauthorized)
		return nil
	}

	o2a := schema.SigninGpaAccountToken{}
	tid := c.Request.FormValue("token_id")
	if tid == "" {
		// tid具有唯一排他索引， 尝试从rid中解析tid
		if idx := strings.IndexRune(rid, '_'); idx > 0 {
			tid = rid[:idx]
		}
	}
	if tid == "" {
		// 兼容方案， 只使用刷新令牌
		if err := o2a.QueryByRefreshToken(a.Sqlx, rid); err != nil {
			if sqlxc.IsNotFound(err) {
				helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID", Other: "令牌无效"}))
			} else {
				helper.FixResponse500Error(c, err, func() {
					logger.Errorf(c, logger.ErrorWW(err))
				})
			}
			return nil
		}
	} else {
		// 使用 TID + RID 方案
		if err := o2a.QueryByTokenKID(a.Sqlx, tid); err != nil {
			if sqlxc.IsNotFound(err) {
				helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID", Other: "令牌无效"}))
			} else {
				helper.FixResponse500Error(c, err, func() {
					logger.Errorf(c, logger.ErrorWW(err))
				})
			}
		} else if o2a.RefreshToken.String != rid {
			// 如果令牌已经被使用
			helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID-2", Other: "令牌已经被消费"}))
		}
	}
	if o2a.ErrCode.Valid && o2a.ErrCode.String != "" {
		// 令牌已经被禁用, 回显令牌禁用原因
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: o2a.ErrCode.String, Other: o2a.ErrMessage.String}))
		return nil
	}
	return &o2a
}

// captcha godoc
// @Tags sign
// @Summary Captcha
// @Description 推送验证码
// @Accept  json
// @Produce  json
// @Param mobile query string false "mobile"
// @Param email query string false "email"
// @Param openid query string false "openid"
// @Param kid query string false "kid"
// @Success 200 {object} helper.Success
// @Router /signin/captcha [get]
func (a *Signin) captcha(c *gin.Context) {
	// 解析参数
	body := schema.SigninOfCaptcha{}
	if err := helper.ParseQuery(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 发送验证码
	code, err := a.SigninService.Captcha(c, &body, "sign")
	if err != nil {
		// 发送验证码失败
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	helper.ResSuccess(c, helper.H{
		"status": "ok",
		"code":   code,
	})
}

//==================================================================================================================
//==================================================================================================================
//==================================================================================================================

// 新建3方令牌
// 该方法不好在于， 签发令牌后， 令牌有可能一次也不会使用， 所以这里应该对令牌进行二次签名
func (a *Signin) tokenNew(c *gin.Context) {

	// 通过刷新令牌生成新令牌
	token, user, err := a.Auther.RefreshToken(c, "", nil)
	// 刷新令牌放生了异常， 直接结束
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	delay := token.GetTokenID() + "_" + crypto.UUID(21)
	a.logSignIn(c, user, token, false, delay)
	// 返回正常结果即可
	helper.ResSuccess(c, helper.H{"token": delay})

}

// 获取3方令牌
func (a *Signin) tokenGet(c *gin.Context) {
	// 需要注意, 刷新令牌只有一次有效
	tid := c.Request.FormValue("token")
	if tid == "" {
		helper.ResError(c, helper.Err401Unauthorized)
		return
	}
	o2a := schema.SigninGpaAccountToken{}
	if err := o2a.QueryByDelayToken(a.Sqlx, tid); err != nil {
		if sqlxc.IsNotFound(err) {
			helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID", Other: "令牌无效"}))
		} else {
			helper.FixResponse500Error(c, err, func() {
				logger.Errorf(c, logger.ErrorWW(err))
			})
		}
	}
	if o2a.DelayExpAt.Int64 < time.Now().Unix() {
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-EXPIRED", Other: "令牌过期"}))
		return
	}
	if o2a.ErrCode.Valid && o2a.ErrCode.String != "" {
		// 令牌已经被禁用, 回显令牌禁用原因
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: o2a.ErrCode.String, Other: o2a.ErrMessage.String}))
		return
	}
	// 令牌结果
	result := schema.SigninResult{
		TokenStatus:  "ok",
		TokenType:    "Bearer",
		TokenID:      o2a.TokenID,
		AccessToken:  o2a.AccessToken.String,
		ExpiresAt:    o2a.ExpiresAt.Int64,
		ExpiresIn:    o2a.ExpiresAt.Int64 - time.Now().Unix(),
		RefreshToken: o2a.RefreshToken.String,
		RefreshExpAt: o2a.RefreshExpAt.Int64,
	}
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
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
