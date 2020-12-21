package middleware

import (
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

var (
	// CasbinNoSignin 未登陆
	CasbinNoSignin = "nosignin"
	// CasbinNoRole 无角色
	CasbinNoRole = "norole"
	// CasbinRolePrefix 角色前缀
	CasbinRolePrefix = "r:"
	// CasbinUserPrefix 角色前缀
	CasbinUserPrefix = "u:"
)

// UserAuthCasbinMiddleware 用户授权中间件
func UserAuthCasbinMiddleware(auther auth.Auther, enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	return UserAuthCasbinMiddlewareByPathFunc(auther, enforcer, func(c *gin.Context, k string) (string, error) {
		return "default", nil
	}, skippers...)
}

// UserAuthCasbinMiddlewareByPathFunc 用户授权中间件
func UserAuthCasbinMiddlewareByPathFunc(auther auth.Auther, enforcer *casbin.SyncedEnforcer, fixOriginFunc func(*gin.Context, string) (string, error), skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return EmptyMiddleware()
	}
	conf := config.C.Casbin

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		user, err := auther.GetUserInfo(c)
		if !conf.Enable {
			// casbin禁用, 只判定是否登陆
			if err != nil {
				// 获取登陆信息异常
				if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
					helper.ResError(c, helper.Err401Unauthorized)
					return
				}
				helper.ResError(c, helper.Err400BadRequest)
				return
			}
			helper.SetUserInfo(c, user)
			c.Next()
			return
		}

		// 需要执行casbin授权
		var sub, usr, aud string      // 角色, 用户,  jwt授权方
		erm := helper.Err403Forbidden // casbin验证失败后返回的异常
		if err != nil {
			if err == auth.ErrNoneToken && conf.NoSignin {
				sub = CasbinNoSignin            // 用户未登陆,且允许执行未登陆认证
				erm = helper.Err401Unauthorized // 替换403异常,因为当前用户未登陆
			} else if err == auth.ErrInvalidToken || err == auth.ErrNoneToken {
				helper.ResError(c, helper.Err401Unauthorized) // 无有效登陆用户
				return
			} else {
				helper.ResError(c, helper.Err500InternalServer) // 解析jwt令牌出现未知错误
				return
			}
		} else {
			ur := user.GetRoleID() // 请求角色
			if ur == "" {
				if conf.NoRole {
					sub = CasbinNoRole // 用户无角色,且允许执行无角色认证,
				} else {
					helper.ResError(c, helper.Err403Forbidden) // 无角色,禁止访问
					return
				}
			} else {
				sub = CasbinRolePrefix + ur
			}
			if !conf.NoUser {
				usr = CasbinUserPrefix + user.GetUserID() // 用户
			}
			aud = user.GetAudience() // jwt授权方
		}
		pat, err := fixOriginFunc(c, helper.XReqOriginPathKey) // 请求路径
		//log.Println(pat)
		if err != nil {
			helper.ResError(c, erm)
			return
		} else if pat == "default" {
			pat = c.Request.URL.Path
		}
		dom, err := fixOriginFunc(c, helper.XReqOriginHostKey) // 请求域名
		//log.Println(dom)
		if err != nil {
			helper.ResError(c, erm)
			return
		} else if dom == "default" {
			dom = c.Request.Host
		}
		act, err := fixOriginFunc(c, helper.XReqOriginMethodKey) // 请求方法
		//log.Println(act)
		if err != nil {
			helper.ResError(c, erm)
			return
		} else if act == "default" {
			act = c.Request.Method
		}
		cip := helper.GetClientIP(c) // 客户端IP
		if b, err := enforcer.Enforce(sub, usr, dom, aud, pat, cip, act); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // 授权发生异常
			helper.ResError(c, erm)
			return
		} else if !b {
			helper.ResError(c, erm)
			return
		}

		if user != nil {
			helper.SetUserInfo(c, user)
		}
		c.Next()
	}
}
