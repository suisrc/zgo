package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/config"
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
func (a *Signin) Register(r gin.IRouter) {
	// sign 开头的路由会被全局casbin放行
	r.POST("signin", a.signin) // 登陆必须是POST请求

	// ua := middleware.UserAuthMiddleware(a.Auther)
	// r.GET("signout", ua, a.signout)

	r.GET("signout", a.signout)
	r.GET("signin/refresh", a.refresh)
	r.GET("signin/captcha", a.captcha)
	//r.GET("signin/mfa", a.signinMFA)

	r.POST("signup", a.signup) // 注册

	r.GET("signin/oauth2/:kid", a.oauth2) // OAUTH2登陆使用GET请求

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
	user, err := a.SigninService.Signin(c, &body, a.last)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	token, usr, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	a.log(c, usr, token, "signin", token.GetRefreshToken())
	// 登陆结果
	result := schema.SigninResult{
		TokenStatus:  "ok",
		TokenType:    "Bearer",
		AccessToken:  token.GetAccessToken(),
		ExpiresAt:    token.GetExpiresAt(),
		ExpiresIn:    token.GetExpiresAt() - time.Now().Unix(),
		RefreshToken: token.GetRefreshToken(),
	}

	// 记录登陆
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}

//==================================================================================================================

// 获取最后一次登陆信息
func (a *Signin) last(c *gin.Context, aid, cid int) (*schema.SigninGpaAccountToken, error) {
	if config.C.JWTAuth.LimitTime <= 0 {
		// 不使用上去签名的结果作为缓存
		return nil, nil
	}
	o2a := schema.SigninGpaAccountToken{}
	if err := o2a.QueryByAccountAndClient(a.Sqlx, aid, cid, helper.GetClientIP(c)); err != nil {
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
			AccessToken:  o2a.AccessToken.String,
			ExpiresAt:    o2a.ExpiresAt.Int64,
			ExpiresIn:    o2a.ExpiresAt.Int64 - time.Now().Unix(),
			RefreshToken: o2a.RefreshToken.String,
		})
	}
	return &o2a, nil
}

// 登陆日志
func (a *Signin) log(c *gin.Context, u auth.UserInfo, t auth.TokenInfo, mode, refresh string) {
	aid, _ := strconv.Atoi(u.GetAccountID())
	cid, cok := helper.GetJwtKidStr(c)
	o2a := schema.SigninGpaAccountToken{
		AccountID:    aid,
		TokenID:      u.GetTokenID(),
		UserKID:      u.GetUserID(),
		RoleKID:      sql.NullString{Valid: true, String: u.GetRoleID()},
		ClientID:     sql.NullInt64{Valid: false},
		ClientKID:    sql.NullString{Valid: cok, String: cid},
		LastIP:       sql.NullString{Valid: true, String: helper.GetClientIP(c)},
		LastAt:       sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:     sql.NullTime{Valid: false},
		LimitKey:     sql.NullString{Valid: false},
		Mode:         sql.NullString{Valid: mode != "", String: mode},
		ExpiresAt:    sql.NullInt64{Valid: t.GetExpiresAt() > 0, Int64: t.GetExpiresAt()},
		AccessToken:  sql.NullString{Valid: t.GetAccessToken() != "", String: t.GetAccessToken()},
		RefreshToken: sql.NullString{Valid: refresh != "", String: refresh},
		Status:       sql.NullBool{Valid: true, Bool: true},
	}
	if _, err := o2a.UpdateAndSaveByAccountAndClient(a.Sqlx); err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
	}

}

//==================================================================================================================

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

	helper.ResSuccess(c, "ok")
}

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
	refreshToken := c.Request.FormValue("refresh_token")
	if refreshToken == "" {
		helper.ResError(c, helper.Err401Unauthorized)
		return
	}
	o2a := schema.SigninGpaAccountToken{}
	if err := o2a.QueryByRefreshToken(a.Sqlx, refreshToken); err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	if o2a.Status.Valid && !o2a.Status.Bool {
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TOKEN-DISABLE", Other: "刷新令牌被禁用"}))
		return
	}

	if config.C.JWTAuth.LimitTime > 0 && o2a.LastAt.Valid && time.Now().Unix()-o2a.LastAt.Time.Unix() < config.C.JWTAuth.LimitTime {
		// 登陆时间非常短,直接返回上次签名结果, 注意, 如果用于在很短时间从两个不同的设备登陆,会导致签发的令牌相同,并且可能会发生同时退出的问题
		// 如果需要避免上述问题,可以禁用缓存
		result := schema.SigninResult{
			TokenStatus: "ok",
			TokenType:   "Bearer",
			AccessToken: o2a.AccessToken.String,
			ExpiresAt:   o2a.ExpiresAt.Int64,
			ExpiresIn:   o2a.ExpiresAt.Int64 - time.Now().Unix(),
		}
		helper.ResSuccess(c, &result)
		return
	}

	token, user, err := a.Auther.RefreshToken(c, o2a.AccessToken.String, func(usr auth.UserInfo, exp int) error {
		if time.Now().Sub(o2a.CreatedAt.Time) > time.Duration(exp)*time.Second {
			// return errors.New("token is expired")
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-TOKEN-EXPIRED", Other: "刷新令牌过期"})
		}
		return nil
	})
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	a.log(c, user, token, "refresh", "")
	result := schema.SigninResult{
		TokenStatus: "ok",
		TokenType:   "Bearer",
		AccessToken: token.GetAccessToken(),
		ExpiresAt:   token.GetExpiresAt(),
		ExpiresIn:   token.GetExpiresAt() - time.Now().Unix(),
		// RefreshToken: token.GetRefreshToken(),
	}
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
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
	code, err := a.SigninService.Captcha(c, &body)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	helper.ResSuccess(c, helper.H{
		"status": "ok",
		"code":   code,
	})
}

// oauth2 godoc
// @Tags sign
// @Summary OAuth2
// @Description 第三方授权登陆
// @Accept  json
// @Produce  json
// @Param kid path string true "kid"
// @Param redirect_uri query string false "redirect_uri"
// @Success 200 {object} helper.Success
// @Router /signin/oauth2/{kid} [get]
func (a *Signin) oauth2(c *gin.Context) {
	// 解析参数
	body := schema.SigninOfOAuth2{}
	if err := helper.ParseQuery(c, &body); err != nil {
		helper.FixResponse406Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 执行登录
	user, err := a.SigninService.OAuth2(c, &body, a.last)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	token, usr, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	a.log(c, usr, token, "oauth2", token.GetRefreshToken())
	// 登陆结果
	result := schema.SigninResult{
		TokenStatus:  "ok",
		TokenType:    "Bearer",
		AccessToken:  token.GetAccessToken(),
		ExpiresAt:    token.GetExpiresAt(),
		ExpiresIn:    token.GetExpiresAt() - time.Now().Unix(),
		RefreshToken: token.GetRefreshToken(),
	}

	// 记录登陆
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}

// Signup godoc
// @Tags sign
// @Summary Signup
// @Description 登陆
// @Accept  json
// @Produce  json
// @Success 200 {object} helper.Success
// @Router /signup [post]
func (a *Signin) signup(c *gin.Context) {
	helper.ResSuccess(c, "功能未开放")
}
