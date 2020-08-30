package api

import (
	"database/sql"
	"strconv"
	"time"

	i18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"
)

// Signin signin
type Signin struct {
	service.GPA
	Auther        auth.Auther
	SigninService service.Signin
}

// Register 注册路由,认证接口特殊,需要独立注册
func (a *Signin) Register(r gin.IRouter) {
	// sign 开头的路由会被全局casbin放行
	r.POST("signin", a.signin)         // 登陆必须是POST请求
	r.POST("signin/{:kid}", a.signin2) // 登陆必须是POST请求

	// ua := middleware.UserAuthMiddleware(a.Auther)
	// r.GET("signout", ua, a.signout)
	r.GET("signout", a.signout)
	r.GET("signin/refresh", a.refresh)

	r.POST("signup", a.signup)         // 注册
	r.POST("signup/{:kid}", a.signup2) // 注册
}

// Signin godoc
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

	// 获取上次登陆的信息
	lastSignin := func(aid, cid int) (*schema.SigninGpaOAuth2Account, error) {
		o2a := schema.SigninGpaOAuth2Account{}
		if err := o2a.QueryByAccountAndClient(a.Sqlx, aid, cid); err != nil {
			if !sqlxc.IsNotFound(err) {
				// 数据库查询发生异常
				logger.Errorf(c, logger.ErrorWW(err))
				return nil, helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-DB-UNKONW", Other: "数据库发生位置异常"})
			}
		}
		if o2a.LastAt.Valid && time.Now().Unix()-o2a.LastAt.Time.Unix() < config.C.JWTAuth.LimitTime {
			// 登陆时间非常短,直接返回上次结果
			return nil, helper.NewSuccess(c, &schema.SigninResult{
				Status:  "ok",
				Token:   o2a.Secret.String,
				Expired: o2a.Expired.Int64,
			})
		}
		return &o2a, nil
	}
	// 执行登录
	user, err := a.SigninService.SigninByPasswd(c, &body, lastSignin)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}
	token, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	// 登陆日志
	aid, _ := strconv.Atoi(user.AccountID)
	cid, cok := helper.GetJwtKidStr(c)
	o2a := schema.SigninGpaOAuth2Account{
		AccountID: aid,
		UserKID:   user.UserID,
		RoleKID:   sql.NullString{Valid: true, String: user.RoleID},
		ClientID:  sql.NullInt64{Valid: false},
		ClientKID: sql.NullString{Valid: cok, String: cid},
		Expired:   sql.NullInt64{Valid: true, Int64: token.GetExpiresAt()},
		LastIP:    sql.NullString{Valid: true, String: helper.GetClientIP(c)},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
		Secret:    sql.NullString{Valid: true, String: token.GetAccessToken()},
		Status:    true,
	}
	if _, err := o2a.UpdateAndSaveByAccountAndClient(a.Sqlx); err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
	}

	// 登陆结果
	result := schema.SigninResult{
		Status:  "ok",
		Token:   token.GetAccessToken(),
		Expired: token.GetExpiresAt(),
		//Expired: token.GetExpiresAt() - time.Now().Unix(),
	}

	// 记录登陆
	// 返回正常结果即可
	helper.ResSuccess(c, &result)
}
func (a *Signin) signin2(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

// Signout godoc
// @Tags sign
// @Summary Signin
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

// Refresh godoc
// @Tags sign
// @Summary Refresh
// @Description 刷新令牌
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} helper.Success
// @Router /signin/refresh [get]
func (a *Signin) refresh(c *gin.Context) {
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
	token, err := a.Auther.GenerateToken(c, user)
	if err != nil {
		helper.FixResponse401Error(c, err, func() {
			logger.Errorf(c, logger.ErrorWW(err))
		})
		return
	}

	result := schema.SigninResult{
		Status:  "ok",
		Token:   token.GetAccessToken(),
		Expired: token.GetExpiresAt(),
		//Expired: token.GetExpiresAt() - time.Now().Unix(),
	}
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
	helper.ResSuccess(c, "功能为开放")
}

// 注册
func (a *Signin) signup2(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}

// 绑定
func (a *Signin) signbind(c *gin.Context) {
	helper.ResSuccess(c, "ok")
}
