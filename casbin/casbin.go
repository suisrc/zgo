package casbin

import (
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/auth"
	"github.com/suisrc/zgo/config"
	"github.com/suisrc/zgo/helper"
	"github.com/suisrc/zgo/logger"
	"github.com/suisrc/zgo/middleware"
	"github.com/suisrc/zgo/store"

	"github.com/gin-gonic/gin"
)

// Implor 外部需要实现的接口
type Implor interface {
	GetAuther() auth.Auther
	GetStorer() store.Storer
	GetTable() string
	GetSqlx2() *sqlx.DB
	GetSuperUserCode() string
	GetPlatformCode() string
	UpdateModelEnable(mid int64) error
	QueryPolicies(org, ver string) (*Policy, error)
	QueryServiceCode(ctx *gin.Context, user auth.UserInfo, host, path, org string) (string, int64, error)
	CheckTenantService(ctx *gin.Context, user auth.UserInfo, org, svc string, sid int64) (bool, error)
}

// 角色定义：
// 1.用户在租户和租户应用上有且各自具有一个角色
// 2.如果在同一个位置(租户或应用)上有多个角色， 服务直接拒绝
// 3.子应用角色优先于租户角色(名称排他除外)
// 3.子应用可以使用使用X-Request-Svc-[SVC-NAME]-Role指定服务角色， 且角色有限被使用

// UseAuthCasbinMiddleware 用户授权中间件
func (a *Auther) UseAuthCasbinMiddleware(skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	return a.UseAuthCasbinMiddlewareByOrigin(func(c *gin.Context, k string) (string, error) { return "default", nil }, skippers...)
}

// UseAuthCasbinMiddlewareByOrigin 用户授权中间件
func (a *Auther) UseAuthCasbinMiddlewareByOrigin(handle func(*gin.Context, string) (string, error), skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return middleware.EmptyMiddleware()
	}
	conf := config.C.Casbin

	return func(c *gin.Context) {
		if middleware.SkipHandler(c, skippers...) {
			c.Next() // 需要跳过权限验证的uri内容
			return
		}

		user, err := a.Implor.GetAuther().GetUserInfo(c, "")
		if err != nil {
			if err == auth.ErrNoneToken || err == auth.ErrInvalidToken {
				helper.ResError(c, helper.Err401Unauthorized)
				return // 无有效登陆用户
			} else if err == auth.ErrExpiredToken {
				helper.ResError(c, helper.Err456TokenExpired)
				return // 访问令牌已经过期
			}
			helper.ResError(c, helper.Err500InternalServer)
			return // 解析jwt令牌出现未知错误
		}
		if !conf.Enable {
			// 禁用了jwt功能
			helper.SetUserInfo(c, user)
			c.Next()
			return
		}
		// 获取访问的域名和路径
		var host, path string // casbin -> 参数
		if host, err = handle(c, helper.XReqOriginHostKey); err != nil {
			helper.ResError(c, helper.Err403Forbidden)
			return
		} else if host == "default" {
			host = c.Request.Host
		}
		if path, err = handle(c, helper.XReqOriginPathKey); err != nil {
			helper.ResError(c, helper.Err403Forbidden)
			return
		} else if path == "default" {
			path = c.Request.URL.Path
		}

		// 获取用户访问的服务
		org := user.GetOrgCode()                                             // casbin -> 参数
		svc, sid, err := a.Implor.QueryServiceCode(c, user, host, path, org) // casbin -> 参数
		if err != nil {
			if err.Error() == "no service" {
				// 访问的服务在权限系统中不存在
				// helper.ResError(c, helper.Err403Forbidden)
				helper.ResError(c, &helper.ErrorModel{
					Status:   403,
					ShowType: helper.ShowWarn,
					ErrorMessage: &i18n.Message{
						ID:    "ERR-SERVICE-NONE",
						Other: "访问的应用不存在",
					},
				})
			} else {
				helper.FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) }) // 未知错误
			}
			return
		}
		c.Writer.Header().Set("X-Request-Z-Svc", svc)
		// 验证服务可访问下
		if b, err := a.Implor.CheckTenantService(c, user, org, svc, sid); err != nil {
			helper.FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) }) // 未知错误
			return
		} else if !b {
			// 租户无法访问该服务
			helper.ResError(c, &helper.ErrorModel{
				Status:   403,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-SERVICE-CLOSE",
					Other: "服务未开通",
				},
			})
			return // 处理过程中发生未知异常
		}
		// 验证用户是否可以跳过权限验证
		if b, err := a.IsPassPermission(c, user, org, svc); err != nil {
			helper.FixResponse403Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) }) // 未知错误
			return
		} else if b {
			// 跳过权限判断
			helper.SetUserInfo(c, user)
			c.Next()
			return
		}

		// 获取用户访问角色
		role, err := a.GetUserRole(c, user, svc, org)
		if err != nil {
			helper.FixResponse403Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) })
			return
		}
		if role == "" {
			helper.ResError(c, &helper.ErrorModel{
				Status:   403,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-SERVICE-NOROLE",
					Other: "用户没有可用角色，拒绝访问",
				},
			})
			return
		}
		c.Writer.Header().Set("X-Request-Z-Svc-Role", role)

		// 租户用户， 默认我们认为租户用户范围不会超过100,000 所以会间人员信息加载到认证器中执行
		// _, _, _ := service.DecryptAccountWithUser(c, user.GetAccount(), user.GetTokenID())
		sub := Subject{
			// UsrID:    aid,
			// AccID:    uid,
			Role:   role,                  // casbin -> 参数 角色
			Acc1:   user.GetAccount1(),    // casbin -> 参数 系统ID
			Acc2:   user.GetAccount2(),    // casbin -> 参数 租户自定义ID
			Usr:    user.GetUserID(),      // casbin -> 参数 用户ID
			Org:    org,                   // casbin -> 参数 租户
			OrgUsr: user.GetOrgUsrID(),    // casbin -> 参数 租户自定义ID
			Iss:    user.GetIssuer(),      // casbin -> 参数 授控域
			Aud:    user.GetAudience(),    // casbin -> 参数 受控域
			Cip:    helper.GetClientIP(c), // casbin -> 参数 ip
			Agent:  user.GetAgent(),       // casbin -> 参数 应用ID
			Scope:  user.GetScope(),       // casbin -> 参数 作用域
		}
		// 访问资源
		method, _ := handle(c, helper.XReqOriginMethodKey)
		obj := Object{
			Svc:    svc,                   // casbin -> 参数 服务
			Host:   host,                  // casbin -> 参数 请求域名
			Path:   path,                  // casbin -> 参数 请求路径
			Method: method,                // casbin -> 参数 请求方法
			Client: helper.GetClientIP(c), // casbin -> 参数 请求IP
		}
		// fix prefix for casbin
		if sub.Usr != "" {
			sub.Usr = UserPrefix + sub.Usr
		}
		if sub.OrgUsr != "" {
			sub.OrgUsr = UserPrefix + sub.OrgUsr
		}
		if sub.Role != "" {
			sub.Role = RolePrefix + sub.Role
		}

		if enforcer, err := a.GetEnforcer(c, user, svc, org); err != nil {
			if helper.FixResponseError(c, err) {
				return
			}
			logger.Errorf(c, logger.ErrorWW(err))
			helper.ResError(c, &helper.ErrorModel{
				Status:   403,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-CASBIN-BUILD",
					Other: "权限验证器发生异常，拒绝访问",
				},
			})
			return
		} else if enforcer == nil {
			// 授权发生异常, 没有可用权限验证器
			helper.ResError(c, helper.Err403Forbidden)
			return
		} else if b, err := enforcer.Enforce(sub, obj); err != nil {
			if helper.FixResponseError(c, err) {
				return
			}
			logger.Errorf(c, logger.ErrorWW(err))
			helper.ResError(c, &helper.ErrorModel{
				Status:   403,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-CASBIN-VERIFY",
					Other: "权限验证器发生异常，拒绝访问",
				},
			})
			return
		} else if !b {
			// 授权失败， 拒绝访问
			// log.Println(ros)
			// log.Println(enforcer.GetImplicitPermissionsForUser(ros))
			helper.ResError(c, helper.Err403Forbidden)
			return
		}

		helper.SetUserInfo(c, user)
		c.Next()
	}
}
