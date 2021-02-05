package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
)

// Signin signin
type Signin struct {
	gpa.GPA
	Auther        auth.Auther
	SigninService *service.Signin
	CasbinAuther  *module.CasbinAuther
}

// Register 注册路由,认证接口特殊,需要独立注册
// sign 开头的路由会被全局casbin放行
func (a *Signin) Register(r gin.IRouter) {

	uax := a.CasbinAuther.UserAuthBasicMiddleware()

	r.POST("signin", a.signin)         // 登录系统， 获取令牌 POST请求
	r.GET("signout", uax, a.signout)   // 登出系统， 注销令牌（访问令牌和刷新令牌）
	r.GET("signin/refresh", a.refresh) // 刷新令牌
	r.GET("signin/captcha", a.captcha) // 发送验证码

	r.POST("pub/3rd/token", a.signin)            // 获取新的访问令牌
	r.GET("pub/3rd/token/new", uax, a.token3new) // 构建新的访问令牌
	r.GET("pub/3rd/token/get", a.token3get)      // 获取新的访问令牌
	r.GET("pub/3rd/token/refresh", a.refresh)    // 获取新的访问令牌

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
	if offset := strings.IndexRune(body.Username, '@'); offset > 0 {
		body.OrgCode = body.Username[offset+1:]
		body.Username = body.Username[:offset]
	}

	// 执行登录， 验证用户
	user, err := a.SigninService.Signin(c, &body, a.LastSignIn)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 执行登录， 生成令牌
	a.signin2(c, user)
}

// 执行登录， 生成令牌
func (a *Signin) signin2(c *gin.Context, u auth.UserInfo) {
	// 生成令牌
	token, usr, err := a.Auther.GenerateToken(c, u)
	if err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 记录登录
	a.LogSignIn(c, usr, token, false, nil)
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
	user, _ := helper.GetUserInfo(c)
	// 执行登出
	if err := a.Auther.DestroyToken(c, user); err != nil {
		helper.ResError(c, helper.Err400BadRequest)
		return
	}
	a.LogSignOut(c, user, user.GetTokenID())

	helper.ResSuccess(c, "ok")
}

//==================================================================================================================

// LastSignIn 获取最后一次登陆信息
func (a *Signin) LastSignIn(c *gin.Context, aid int64) (*schema.SigninGpaAccountToken, error) {
	if config.C.JWTAuth.LimitTime <= 0 {
		// 不使用上去签名的结果作为缓存
		return nil, nil
	}
	o2a := schema.SigninGpaAccountToken{}
	// 防止意外放生， 使用客户端IP作为影响因子
	if err := o2a.QueryByAccountAndClient2(a.Sqlx2, aid, helper.GetClientIP(c)); err != nil {
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
			ExpiresAt:    o2a.ExpiresAt.Time.Unix(),
			ExpiresIn:    int64(o2a.ExpiresAt.Time.Unix() - time.Now().Unix()),
			RefreshToken: o2a.RefreshToken.String,
			RefreshExpAt: o2a.RefreshExpAt.Time.Unix(),
		})
	}
	return &o2a, nil
}

// LogSignIn 日志记录
func (a *Signin) LogSignIn(c *gin.Context, u auth.UserInfo, t auth.TokenInfo, update bool, fix func(*schema.SigninGpaAccountToken)) {
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
		TokenPID:     sql.NullString{Valid: u.GetTokenPID() != "", String: u.GetTokenPID()},
		OrgCode:      sql.NullString{Valid: u.GetOrgCode() != "", String: u.GetOrgCode()},
		AccessToken:  sql.NullString{Valid: t.GetAccessToken() != "", String: t.GetAccessToken()},
		ExpiresAt:    sql.NullTime{Valid: t.GetExpiresAt() > 0, Time: time.Unix(t.GetExpiresAt(), 0)},
		RefreshToken: sql.NullString{Valid: t.GetRefreshToken() != "", String: t.GetRefreshToken()},
		RefreshExpAt: sql.NullTime{Valid: t.GetRefreshExpAt() > 0, Time: time.Unix(t.GetRefreshExpAt(), 0)},
		LastIP:       sql.NullString{Valid: true, String: helper.GetClientIP(c)},
		LastAt:       sql.NullTime{Valid: true, Time: time.Now()},
	}
	if fix != nil {
		fix(&o2a)
	}
	if err := o2a.UpdateAndSaveByTokenKID2(a.Sqlx2, update); err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
	}
}

// LogSignOut 日志记录
func (a *Signin) LogSignOut(c *gin.Context, u auth.UserInfo, t string) {
	// 销毁刷新令牌
	o2a := schema.SigninGpaAccountToken{
		TokenID:      u.GetTokenID(),
		RefreshExpAt: sql.NullTime{Valid: true, Time: time.Unix(1, 0)},
	}
	if err := o2a.UpdateAndSaveByTokenKID2(a.Sqlx2, true); err != nil {
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
	o2a := a.getSigninGpaAccountTokenByRefresh(c)
	if o2a == nil {
		return // 结束处理
	}
	if config.C.JWTAuth.LimitTime > 0 && o2a.LastAt.Valid && time.Now().Unix()-o2a.LastAt.Time.Unix() < config.C.JWTAuth.LimitTime {
		// 登陆时间非常短,直接返回上次签名结果, 注意, 如果用于在很短时间从两个不同的设备登陆,会导致签发的令牌相同,并且可能会发生同时退出的问题
		// 如果需要避免上述问题,可以禁用缓存
		result := schema.SigninResult{
			TokenStatus:  "ok",
			TokenType:    "Bearer",
			TokenID:      o2a.TokenID,
			AccessToken:  o2a.AccessToken.String,
			ExpiresAt:    o2a.ExpiresAt.Time.Unix(),
			ExpiresIn:    o2a.ExpiresAt.Time.Unix() - time.Now().Unix(),
			RefreshToken: o2a.RefreshToken.String,
			RefreshExpAt: o2a.RefreshExpAt.Time.Unix(),
		}
		helper.ResSuccess(c, &result)
		return
	}
	// 通过刷新令牌生成新令牌
	token, user, err := a.Auther.RefreshToken(c, o2a.AccessToken.String, func(usrInfo auth.UserInfo, expIn int) error {
		if o2a.RefreshExpAt.Time.Unix() == 1 {
			// 刷新令牌被销毁
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-EESTROY", Other: "刷新令牌已销毁"})
		} else if o2a.RefreshExpAt.Time.Before(time.Now()) {
			// 刷新令牌已经过期， 无法执行刷新
			// return errors.New("token is expired")
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-EXPIRED", Other: "刷新令牌已过期"})
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
	a.LogSignIn(c, user, token, true, nil)
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
func (a *Signin) getSigninGpaAccountTokenByRefresh(c *gin.Context) *schema.SigninGpaAccountToken {
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
		if err := o2a.QueryByRefreshToken2(a.Sqlx2, rid); err != nil {
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
		if err := o2a.QueryByTokenKID2(a.Sqlx2, tid); err != nil {
			if sqlxc.IsNotFound(err) {
				helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID", Other: "令牌无效"}))
			} else {
				helper.FixResponse500Error(c, err, func() {
					logger.Errorf(c, logger.ErrorWW(err))
				})
			}
			return nil
		} else if o2a.RefreshToken.String != rid {
			// 如果令牌已经被使用
			helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-USED", Other: "令牌已被使用"}))
			return nil
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
func (a *Signin) token3new(c *gin.Context) {
	// 确定登陆用户的身份
	usr, _ := helper.GetUserInfo(c)

	aid, uid, _ := service.DecryptAccountWithUser(c, usr.GetAccount(), usr.GetTokenID())
	tid := jwt.NewTokenID(strconv.Itoa(int(aid)))
	tkn := tid + "_" + crypto.UUID(21)

	// a.logSignIn(c
	o2a := schema.SigninGpaAccountToken{
		TokenID:   tid,
		AccountID: aid,
		OrgCode:   sql.NullString{Valid: usr.GetOrgCode() != "", String: usr.GetOrgCode()},
		Number1:   sql.NullInt64{Valid: true, Int64: uid},
		String1:   sql.NullString{Valid: true, String: usr.GetTokenID()},
		CodeToken: sql.NullString{Valid: true, String: tkn},
		CodeExpAt: sql.NullTime{Valid: true, Time: time.Now().Add(300 * time.Second)},
	}
	if err := o2a.UpdateAndSaveByTokenKID2(a.Sqlx2, false); err != nil {
		helper.FixResponse500Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	// 返回正常结果即可
	helper.ResSuccess(c, helper.H{"status": "ok", "code": tkn})

}

// 获取3方令牌
func (a *Signin) token3get(c *gin.Context) {
	o2a := a.getSigninGpaAccountTokenByCode(c)
	if o2a == nil {
		return // 结束处理
	}
	// 通过刷新令牌生成新令牌
	token, user, err := a.Auther.RefreshToken(c, o2a.AccessToken.String, func(usrInfo auth.UserInfo, expIn int) error {
		if usr, b := usrInfo.(*jwt.UserClaims); b {
			// 修正数据
			usr.Id = o2a.TokenID
			usr.Account, _ = service.EncryptAccountWithUser(c, o2a.AccountID, o2a.Number1.Int64, o2a.TokenID)
			usr.Audience = c.Request.Host
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

	a.LogSignIn(c, user, token, true, func(sga *schema.SigninGpaAccountToken) {
		// 注销延迟令牌， 延迟令牌只允许使用一次
		sga.CodeExpAt = sql.NullTime{Valid: true, Time: time.Unix(1, 0)}
	})
	// 令牌结果
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
func (a *Signin) getSigninGpaAccountTokenByCode(c *gin.Context) *schema.SigninGpaAccountToken {
	// 需要注意, 刷新令牌只有一次有效
	tid := c.Request.FormValue("code")
	if tid == "" {
		helper.ResError(c, helper.Err401Unauthorized)
		return nil
	}
	o2a := schema.SigninGpaAccountToken{}
	if err := o2a.QueryByDelayToken2(a.Sqlx2, tid); err != nil {
		if sqlxc.IsNotFound(err) {
			helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID", Other: "令牌无效"}))
		} else {
			helper.FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) })
		}
		return nil
	}
	if o2a.CodeExpAt.Time.Unix() == 1 {
		// Code令牌被使用
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-USED", Other: "令牌已使用"}))
		return nil
	} else if o2a.CodeExpAt.Time.Before(time.Now()) {
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-EXPIRED", Other: "令牌已过期"}))
		return nil
	}
	if o2a.AccessToken.Valid {
		result := schema.SigninResult{
			TokenStatus:  "ok",
			TokenType:    "Bearer",
			TokenID:      o2a.TokenID,
			AccessToken:  o2a.AccessToken.String,
			ExpiresAt:    o2a.ExpiresAt.Time.Unix(),
			ExpiresIn:    o2a.ExpiresAt.Time.Unix() - time.Now().Unix(),
			RefreshToken: o2a.RefreshToken.String,
			RefreshExpAt: o2a.RefreshExpAt.Time.Unix(),
		}
		helper.ResSuccess(c, &result)
		return nil
	}
	o2b := schema.SigninGpaAccountToken{}
	if err := o2b.QueryByTokenKID2(a.Sqlx2, o2a.String1.String); err != nil {
		if sqlxc.IsNotFound(err) {
			helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-TOKEN-INVALID2", Other: "原始令牌无效"}))
		} else {
			helper.FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) })
		}
		return nil
	}
	if o2b.ErrCode.Valid && o2b.ErrCode.String != "" {
		// 令牌已经被禁用, 回显令牌禁用原因
		helper.ResJSON(c, http.StatusOK, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: o2a.ErrCode.String, Other: o2a.ErrMessage.String}))
		return nil
	}
	o2a.AccessToken = o2b.AccessToken
	return &o2a
}
