package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/casbin/casbin"
	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/middleware"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"

	"github.com/gin-gonic/gin"
)

// CasbinPolicyModel casbin使用的对比模型
var CasbinPolicyModel = `
[request_definition]
r = sub, ros

[policy_definition]
p = sub, dom, pat, cip, act, eft

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = 
`

// (g(r.sub, p.sub) || r.usr != "" && g2(r.usr, p.sub)) && fact(r.act, p.act) && fdoma(r.dom, p.dom, r.aud) && (p.pat=="" || keyMatch(r.pat, p.pat)) && (p.cip=="" || ipMatch(r.cip, p.cip))

var (
	// CasbinNoSign 未登陆
	CasbinNoSign = "nosign"
	// CasbinNoRole 无角色
	CasbinNoRole = "norole"
	// CasbinNoUser 无用户
	CasbinNoUser = "nouser"
	// CasbinRolePrefix 角色
	CasbinRolePrefix = "r:"
	// CasbinUserPrefix 用户
	CasbinUserPrefix = "u:"
	// CasbinPolicyPrefix 策略
	CasbinPolicyPrefix = "p:"
)

// CasbinAuther 权限管理
type CasbinAuther struct {
	gpa.GPA              // 数据库
	Auther  auth.Auther  // 权限
	Storer  store.Storer // 缓存
}

// UserAuthCasbinMiddleware 用户授权中间件
func (a *CasbinAuther) UserAuthCasbinMiddleware(skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	return a.UserAuthCasbinMiddlewareByOrigin(func(c *gin.Context, k string) (string, error) { return "default", nil }, skippers...)
}

// UserAuthCasbinMiddlewareByOrigin 用户授权中间件
func (a *CasbinAuther) UserAuthCasbinMiddlewareByOrigin(handle func(*gin.Context, string) (string, error), skippers ...middleware.SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return middleware.EmptyMiddleware()
	}
	conf := config.C.Casbin

	return func(c *gin.Context) {
		if middleware.SkipHandler(c, skippers...) {
			c.Next() // 需要跳过权限验证的uri内容
			return
		}

		user, err := a.Auther.GetUserInfo(c)
		erm := helper.Err403Forbidden // casbin验证失败后返回的异常
		if err != nil {
			if err == auth.ErrNoneToken && conf.NoSign {
				erm = helper.Err401Unauthorized // 替换403异常,因为当前用户未登陆
			} else if err == auth.ErrNoneToken || err == auth.ErrInvalidToken {
				helper.ResError(c, helper.Err401Unauthorized)
				return // 无有效登陆用户
			} else {
				helper.ResError(c, helper.Err500InternalServer)
				return // 解析jwt令牌出现未知错误
			}
		}
		if !conf.Enable {
			// casbin禁用, 只判定是否登陆
			if err != nil {
				helper.ResError(c, erm)
				return
			}
			helper.SetUserInfo(c, user)
			c.Next()
			return
		}

		host, err := handle(c, helper.XReqOriginHostKey) // casbin -> 参数
		if err != nil {
			helper.ResError(c, erm)
			return
		} else if host == "default" {
			host = c.Request.Host
		}
		path, err := handle(c, helper.XReqOriginPathKey) // casbin -> 参数
		if err != nil {
			helper.ResError(c, erm)
			return
		} else if path == "default" {
			path = c.Request.URL.Path
		}
		org := user.GetOrgCode()                                      // casbin -> 参数
		svc, sid, err := a.QueryServiceCode(c, user, host, path, org) // casbin -> 参数
		if err != nil {
			helper.ResError(c, helper.Err500InternalServer)
			return // 处理过程中发生未知异常
		}
		if b, err := a.CheckTenantService(c, user, org, svc, sid); err != nil {
			helper.FixResponse500Error(c, err, func() {
				logger.Errorf(c, logger.ErrorWW(err))
			})
			return // 处理过程中发生未知异常
		} else if !b {
			// 租户无法访问该服务
			helper.ResError(c, &helper.ErrorModel{
				Status:   200,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-SERVICE-CLOSE",
					Other: "服务未开通",
				},
			})
			return // 处理过程中发生未知异常
		}
		// 确定管理员身份， 这里是否担心管理员身份被篡改？
		// 如果签名密钥泄漏， 会发生签名篡改问题， 所以需要保密服务器签名密钥
		if user.GetOrgAdmin() == schema.SuperUser {
			helper.SetUserInfo(c, user)
			c.Next()
			return // 管理员， 跳过所有验证， 可以访问服务的所有权限
		}

		aid, uid, _ := DecryptAccountWithUser(c, user.GetAccount(), user.GetTokenID())
		method, _ := handle(c, helper.XReqOriginMethodKey)
		sub := CasbinSubject{
			UsrID:    aid,
			AccID:    uid,
			Host:     host,
			Path:     path,
			Org:      org,
			Svc:      svc,
			Iss:      user.GetIssuer(),          // casbin -> 参数 授控域
			Aud:      user.GetAudience(),        // casbin -> 参数 受控域
			Usr:      user.GetUserID(),          // casbin -> 参数 用户ID
			OrgUsr:   user.GetOrgUsrID(),        // casbin -> 参数 租户自定义ID
			OrgApp:   user.GetOrgAppID(),        // casbin -> 参数 应用ID
			Roles:    user.GetUserRoles(),       // casbin -> 参数 角色
			RolesSvc: user.GetUserSvcRoles(svc), // casbin -> 参数 角色
			ClientIP: helper.GetClientIP(c),     // casbin -> 参数 请求IP
			Method:   method,                    // casbin -> 参数 请求方法
		}

		if enforcer, err := a.GetEnforcer(c, user, svc, org); err != nil {
			logger.Errorf(c, logger.ErrorWW(err))
			helper.ResError(c, erm) // 授权发生异常
			return
		} else if enforcer == nil {
			helper.ResError(c, erm) // 授权发生异常
			return
		} else if b := enforcer.Enforce(sub); !b {
			helper.ResError(c, erm) // 授权失败
			return
		}

		if user != nil {
			helper.SetUserInfo(c, user)
		}
		c.Next()
	}
}

//======================================================================================
//======================================================================================
//======================================================================================

// GetEnforcer 获取验证控制器
func (a *CasbinAuther) GetEnforcer(ctx *gin.Context, user auth.UserInfo, svc, org string) (*casbin.SyncedEnforcer, error) {
	return nil, errors.New("no casbin")
}

// CasbinSubject subject
type CasbinSubject struct {
	UsrID    int
	AccID    int
	Host     string
	Path     string
	Org      string
	Svc      string
	Iss      string
	Aud      string
	Usr      string
	OrgUsr   string
	OrgApp   string
	Roles    []string
	RolesSvc []string
	ClientIP string
	Method   string
}

// CasbinPolicy casbin策略
type CasbinPolicy struct {
	ModelText   string   // 模型声明
	PolicyLines []string // 策略声明
}

// QueryCasbinPolicy 获取Casbin策略
func QueryCasbinPolicy(sqlx *sqlx.DB, org string) (*CasbinPolicy, error) {
	return nil, nil
}

// QueryServiceCode 查询服务
// "fmes:svc-code:" + host + ":" + resource
func (a *CasbinAuther) QueryServiceCode(ctx *gin.Context, user auth.UserInfo, host, path, org string) (string, int, error) {
	rescount := 3
	offset := strings.IndexFunc(path, func(r rune) bool {
		if r == '/' {
			if rescount--; rescount == 0 {
				return true
			}
		}
		return false
	})
	resource := ""
	if offset > 0 && strings.HasPrefix(path[:offset], "/api/") {
		resource = path[:offset]
	}
	key := "fmes:svc-code:" + host + ":" + resource

	if svc, b, err := a.Storer.Get(ctx, key); err != nil {
		return "", 0, err
	} else if b {
		if svc == "err" {
			// 上一次查询， 执行了拒绝请求
			return "", 0, errors.New("no service")
		}
		offset := strings.IndexRune(svc, '/')
		sid, _ := strconv.Atoi(svc[offset+1:])
		return svc[:offset], sid, nil
	}

	sa := schema.CasbinGpaSvcAud{}
	if err := sa.QueryByAudAndResAndOrg(a.Sqlx, host, resource, ""); err != nil {
		a.Storer.Set(ctx, key, "err", 60*time.Second) // 1分钟延迟刷新， 拒绝请求也需要缓存
		return "", 0, err
	} else if sa.SvcCode.Valid {
		a.Storer.Set(ctx, key, "err", 60*time.Second) // 1分钟延迟刷新， 拒绝请求也需要缓存
		return "", 0, errors.New("no service")
	}
	svc := sa.SvcCode.String + "/" + strconv.Itoa(int(sa.SvcID.Int64))
	a.Storer.Set(ctx, key, svc, 300*time.Second) // 查询结果缓存5分钟
	return sa.SvcCode.String, int(sa.SvcID.Int64), nil
}

// CheckTenantService 验证租户是否有访问该服务的权限服务
// "fmes:svc-org:" + svc_cod + ":" + org_cod -> CasbinGpaSvcOrg
func (a *CasbinAuther) CheckTenantService(ctx *gin.Context, user auth.UserInfo, org, svc string, sid int) (bool, error) {
	if org == "" || org == schema.PlatformCode {
		return true, nil // 平台用户， 没有服务权限问题
	}
	key := "fmes:svc-org:" + svc + ":" + org

	if res, b, err := a.Storer.Get(ctx, key); err != nil {
		return false, err
	} else if b {
		return res == "1", nil
	}

	var erm error = nil
	so := schema.CasbinGpaSvcOrg{}
	// 1:启用 0:禁用 2:未激活 3: 删除 4: 欠费 5: 到期
	if err := so.QueryByOrgAndSvc(a.Sqlx, org, sid); err != nil {
		erm = err
	} else if so.Expired.Valid && time.Now().After(so.Expired.Time) {
		erm = helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"})
	} else if so.Status == schema.StatusEnable {
		// 正常结果返回
		expiration := 300 * time.Second // 5分钟延迟刷新
		if so.Expired.Valid && so.Expired.Time.Sub(time.Now()) < expiration {
			// 修正过期时间
			expiration = so.Expired.Time.Sub(time.Now())
		}
		a.Storer.Set(ctx, key, "1", expiration) // 查询结果缓存5分钟
		return true, nil
	} else if so.Status == schema.StatusDisable {
		erm = helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: "WARN-SERVICE-DISABLE", Other: "服务已经被禁用"})
	} else if so.Status == schema.StatusDeleted {
		erm = helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: "WARN-SERVICE-DELETE", Other: "服务已经被删除"})
	} else if so.Status == schema.StatusNoActivate {
		erm = helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: "WARN-SERVICE-NOACTIVATE", Other: "服务未激活"})
	} else if so.Status == schema.StatusExpired {
		erm = helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"})
	}
	a.Storer.Set(ctx, key, "0", 60*time.Second) // 1分钟延迟刷新， 拒绝请求也需要缓存
	return false, erm
}
