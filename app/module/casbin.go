package module

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/middleware"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"

	"github.com/gin-gonic/gin"

	zgocasbin "github.com/suisrc/zgo/modules/casbin"
)

// 角色定义：
// 1.用户在租户和租户应用上有且各自具有一个角色
// 2.如果在同一个位置(租户或应用)上有多个角色， 服务直接拒绝
// 3.子应用角色优先于租户角色(名称排他除外)
// 3.子应用可以使用使用X-Request-Svc-[SVC-NAME]-Role指定服务角色， 且角色有限被使用

// CasbinAuther 权限管理
type CasbinAuther struct {
	gpa.GPA                                   // 数据库
	Auther         auth.Auther                // 权限
	Storer         store.Storer               // 缓存
	CachedEnforcer map[string]*CasbinEnforcer // 验证器
	CachedExpireAt time.Time                  // 刷新时间
	Mutex          sync.RWMutex               // 同步锁
}

// CasbinEnforcer 验证器
type CasbinEnforcer struct {
	Enforcer *casbin.SyncedEnforcer // 验证器
	ExpireAt time.Time              // 过期时间
	Refresh  time.Time              // 刷新时间
	Version  string                 // 验证版本
	Mutex    sync.RWMutex           // 同步锁
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
		if err != nil {
			if err == auth.ErrNoneToken || err == auth.ErrInvalidToken {
				helper.ResError(c, helper.Err401Unauthorized)
				return // 无有效登陆用户
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
		org := user.GetOrgCode()                                      // casbin -> 参数
		svc, sid, err := a.QueryServiceCode(c, user, host, path, org) // casbin -> 参数
		if err != nil {
			helper.ResError(c, helper.Err500InternalServer)
			return // 处理过程中发生未知异常
		}
		// 验证服务可访问下
		if b, err := a.CheckTenantService(c, user, org, svc, sid); err != nil {
			helper.FixResponse500Error(c, err, func() {
				logger.Errorf(c, logger.ErrorWW(err))
			}) // 未知错误
			return
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
		// 验证用户是否可以跳过权限验证
		if b, err := a.IsPassPermission(c, user, org, svc); err != nil {
			helper.FixResponse403Error(c, err, func() {
				logger.Errorf(c, logger.ErrorWW(err))
			}) // 未知错误
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
			helper.FixResponse403Error(c, err, func() {
				logger.Errorf(c, logger.ErrorWW(err))
			})
			return
		}
		if role == "" {
			helper.ResError(c, &helper.ErrorModel{
				Status:   200,
				ShowType: helper.ShowWarn,
				ErrorMessage: &i18n.Message{
					ID:    "ERR-SERVICE-NOROLE",
					Other: "用户没有可用角色， 拒绝访问",
				},
			})
			return
		}

		// 租户用户， 默认我们认为租户用户范围不会超过100,000 所以会间人员信息加载到认证器中执行
		// aid, uid, cus, _ := DecryptAccountWithUser(c, user.GetAccount(), user.GetTokenID())
		sub := CasbinSubject{
			// UsrID:    aid,
			// AccID:    uid,
			Org:    org,                // casbin -> 参数 租户
			Iss:    user.GetIssuer(),   // casbin -> 参数 授控域
			Aud:    user.GetAudience(), // casbin -> 参数 受控域
			Usr:    user.GetUserID(),   // casbin -> 参数 用户ID
			OrgUsr: user.GetOrgUsrID(), // casbin -> 参数 租户自定义ID
			OrgApp: user.GetOrgAppID(), // casbin -> 参数 应用ID
		}
		// 访问资源
		method, _ := handle(c, helper.XReqOriginMethodKey)
		obj := CasbinObject{
			Svc:    svc,                   // casbin -> 参数 服务
			Host:   host,                  // casbin -> 参数 请求域名
			Path:   path,                  // casbin -> 参数 请求路径
			Method: method,                // casbin -> 参数 请求方法
			Client: helper.GetClientIP(c), // casbin -> 参数 请求IP
		}
		// fix prefix for casbin
		if sub.Usr != "" {
			sub.Usr = CasbinUserPrefix + sub.Usr
		}
		if sub.OrgUsr != "" {
			sub.OrgUsr = CasbinUserPrefix + sub.OrgUsr
		}
		ros := CasbinRolePrefix + role

		if enforcer, err := a.GetEnforcer(conf, c, user, svc, org); err != nil {
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
		} else if b, err := enforcer.Enforce(sub, obj, ros); err != nil {
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

		c.Writer.Header().Set("X-Request-Z-Svc", svc)
		c.Writer.Header().Set("X-Request-Z-Svc-Role", role)
		helper.SetUserInfo(c, user)
		c.Next()
	}
}

//======================================================================================
//======================================================================================
//======================================================================================

// IsPassPermission 跳过权限判断
// 确定管理员身份， 这里是否担心管理员身份被篡改？如果签名密钥泄漏， 会发生签名篡改问题， 所以需要保密服务器签名密钥
func (a *CasbinAuther) IsPassPermission(c *gin.Context, user auth.UserInfo, svc, org string) (bool, error) {
	if user.GetOrgAdmin() == schema.SuperUser {
		// 组织管理员， 跳过验证
		return true, nil
	} else if user.GetOrgCode() == schema.PlatformCode {
		// 平台用户， 暂不处理
	} else if user.GetOrgCode() == "" {
		// 无租户用户, 只验证登录， 注意：无组织用户只能访问pub服务
		if strings.HasPrefix(svc, CasbinSvcPublic) {
			return true, nil
		}
		return false, &helper.ErrorModel{
			Status:   200,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-SERVICE-TENANT-NONE",
				Other: "无租户信息，拒绝访问",
			},
		}
	}
	return false, nil
}

// GetUserRole 获取验证控制器
func (a *CasbinAuther) GetUserRole(c *gin.Context, user auth.UserInfo, svc, org string) (role string, err error) {
	if roles := user.GetUserRoles(); len(roles) == 0 {
		// 当前用户没有可用角色
		return "", nil
	} else if len(roles) == 1 {
		// 当前用户只有一个角色
		return roles[0], nil
	}
	// 处理多角色问题
	roles := []string{}
	if svc != "" && svc != schema.PlatformCode {
		role = c.GetHeader(strings.ReplaceAll(CasbinSvcRoleKey, "[SVC-NAME]", svc)) // 子应用， 需要子应用授权
		if role != "" {
			// 验证角色信息， 快速结束
			for _, v := range user.GetUserRoles() {
				if role == v {
					return role, nil // 角色有效直接返回
				}
			}
		}
		// 无指定角色， 获取用户服务角色
		roles = user.GetUserSvcRoles(svc) // 有可能没有角色信息， 使用len(roles) == 0判断没有应用角色
	}
	if role == "" && len(roles) == 0 {
		role = c.GetHeader(CasbinSysRoleKey) // 使用系统平台角色
		if role != "" {
			// 验证角色信息， 快速结束
			for _, v := range user.GetUserRoles() {
				if role == v {
					return role, nil // 角色有效直接返回
				}
			}
		}
		// 无指定角色， 获取用户服务角色
		roles = user.GetUserRoles()
	}
	if role != "" {
		// 指定的角色无效
		err = &helper.ErrorModel{
			Status:   200,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-SERVICE-ROLE-INVALID",
				Other: "用户指定的角色无效",
			},
		}
	} else if len(roles) == 1 {
		role = roles[0] // 只有单角色， 配置用户角色
	} else if len(roles) > 1 {
		// 用户对同一个应用具有多个角色， 拒绝访问
		err = &helper.ErrorModel{
			Status:   200,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-SERVICE-ROLE-MULT",
				Other: "用户访问的应用同时具有多角色，且没有指定角色",
			},
		}
	}
	return
}

// ClearEnforcer 清理缓存
func (a *CasbinAuther) ClearEnforcer(force bool, org string) {
	if a.CachedEnforcer == nil {
		// do nothing
	} else if force {
		a.CachedEnforcer = map[string]*CasbinEnforcer{} // 删除之前所有的
		a.CachedExpireAt = time.Now().Add(time.Minute)
	} else if org != "" {
		key := "fmes:casbin:" + org
		delete(a.CachedEnforcer, key) // 清除指定缓存
	} else {
		now := time.Now()
		for k, v := range a.CachedEnforcer {
			if v.ExpireAt.Before(now) {
				delete(a.CachedEnforcer, k) // 清除过期
			}
		}
	}
}

// GetEnforcer 获取验证控制器
func (a *CasbinAuther) GetEnforcer(conf config.Casbin, ctx *gin.Context, user auth.UserInfo, svc, org string) (*casbin.SyncedEnforcer, error) {
	if a.CachedEnforcer == nil {
		a.CachedEnforcer = map[string]*CasbinEnforcer{}
		a.CachedExpireAt = time.Now().Add(2 * time.Minute)
	} else if a.CachedExpireAt.Before(time.Now()) {
		a.CachedExpireAt = time.Now().Add(2 * time.Minute) // 设定1分钟后再次刷新
		go a.ClearEnforcer(false, "")                      // 执行异步刷新流程
	}
	key := "fmes:casbin:" + org
	ver := ""
	cached, exist := a.CachedEnforcer[key]
	if exist && cached.Refresh.Before(time.Now()) {
		return cached.Enforcer, nil
	} else if exist {
		ver = cached.Version
		// 多进程更新
		defer cached.Mutex.Unlock()
		cached.Mutex.Lock()
		if c2, e2 := a.CachedEnforcer[key]; e2 && c2.Refresh.Before(time.Now()) {
			return c2.Enforcer, nil // 缓存已经由其他进程更新
		}
	} else {
		// 多进程创建
		defer a.Mutex.Unlock()
		a.Mutex.Lock()
		if c2, e2 := a.CachedEnforcer[key]; e2 {
			return c2.Enforcer, nil // 缓存已经由其他进程更新
		}

	}
	c, err := a.QueryCasbinPolicy(org, ver)
	if err != nil {
		return nil, err
	} else if c == nil {
		// 版本不变, 重置有效期限
		cached.ExpireAt = time.Now().Add(10 * time.Minute)
		cached.Refresh = time.Now().Add(time.Minute)
		return cached.Enforcer, nil
	}
	if c.ModelText == "" {
		// 重新加载配置, *Adapter
		// 数据库访问适配器（使用redis缓存请改写这里）
		if adapter, b := cached.Enforcer.GetAdapter().(*Adapter); b {
			if adapter.Mid != c.Mid || adapter.Ver != c.Ver {
				adapter.Mid = c.Mid
				adapter.Ver = c.Ver
			}
		} else {
			return nil, errors.New("casbin adapter type is error")
		}
		cached.Enforcer.LoadPolicy()
		cached.ExpireAt = time.Now().Add(10 * time.Minute)
		cached.Refresh = time.Now().Add(time.Minute)
		cached.Version = c.Version
		return cached.Enforcer, nil
	}

	// log.Println(c)
	m, err := model.NewModelFromString(c.ModelText)
	if err != nil {
		return nil, err
	}
	// *Adapter
	// 数据库访问适配器（使用redis缓存请改写这里）
	adapter := NewCasbinAdapter(a.Sqlx, schema.TableCasbinRule, c.Mid, c.Ver)
	e, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	e.EnableLog(conf.Debug)
	e.EnableEnforce(conf.Enable)

	// 增加策略关系
	e.AddNamedPolicies("p", c.Policies)
	// 增加策略关系
	e.AddNamedGroupingPolicies("g", c.Grouping)

	// 注册方法
	e.AddFunction("domainMatch", zgocasbin.DomainMatchFunc)
	e.AddFunction("methodMatch", zgocasbin.MethodMatchFunc)
	e.AddFunction("audienceMatch", zgocasbin.AudienceMatchFunc)

	cgm := schema.CasbinGpaModel{
		ID:     c.Mid,
		Status: schema.StatusEnable,
	}
	// 更新状态
	if err := cgm.SaveOrUpdate(a.Sqlx); err != nil {
		return nil, err
	}

	// 配置缓存
	a.CachedEnforcer[key] = &CasbinEnforcer{
		Enforcer: e,
		ExpireAt: time.Now().Add(10 * time.Minute), // 有效期10分钟
		Refresh:  time.Now().Add(time.Minute),      // 刷新时间1分钟
		Version:  ver,
	}
	return e, nil
}

// CasbinPolicy Casbin策略
type CasbinPolicy struct {
	Mid       int64
	Ver       string
	ModelText string     // 模型声明
	Grouping  [][]string // 角色声明
	Policies  [][]string // 策略声明
	Version   string     // 策略版本
}

// QueryCasbinPolicy 获取Casbin策略
func (a *CasbinAuther) QueryCasbinPolicy(org, ver string) (*CasbinPolicy, error) {
	c := CasbinPolicy{
		Grouping: [][]string{},
		Policies: [][]string{},
	}
	// 获取策略模型
	cgm := schema.CasbinGpaModel{}
	if err := cgm.QueryByOrg(a.Sqlx, org); !sqlxc.IsNotFound(err) {
		// 数据库异常
		return nil, err
	} else if ver != "" && cgm.ID > 0 && ver == fmt.Sprintf("%s:%s", strconv.Itoa(int(cgm.ID)), cgm.UseVer.String) {
		// 访问策略不变
		return nil, nil
	} else if cgm.ID > 0 {
		// 访问策略更新
		c.Mid = cgm.ID
		c.Ver = cgm.UseVer.String
		c.Version = fmt.Sprintf("%s:%s", strconv.Itoa(int(cgm.ID)), cgm.UseVer.String)
		if cgm.Status == schema.StatusEnable {
			return &c, nil
		}
	}
	// 获取基础配置访问策略
	if err := a.CreateCasbinPolicy(org, &c); err != nil {
		return nil, err
	}
	// 待激活， 完成激活操作
	if cgm.Statement.Valid {
		c.ModelText = CasbinPolicyModel + cgm.Statement.String
	} else {
		c.ModelText = CasbinPolicyModel + CasbinDefaultMatcher
	}
	if cgm.ID > 0 {
		return &c, nil
	}

	// 新建访问策略
	cgm.Name = "Default"
	cgm.UseVer = sql.NullString{Valid: true, String: "1.0.0"}
	cgm.OrgCode = sql.NullString{Valid: true, String: org}
	cgm.Statement = sql.NullString{Valid: true, String: CasbinDefaultMatcher}
	cgm.Description = sql.NullString{Valid: true, String: "Auto Build"}
	cgm.Status = schema.StatusNoActivate // 未激活状态
	if err := cgm.SaveOrUpdate(a.Sqlx); err != nil {
		return nil, err
	}
	c.Mid = cgm.ID
	c.Ver = cgm.UseVer.String
	c.Version = fmt.Sprintf("%s:%s", strconv.Itoa(int(cgm.ID)), cgm.UseVer.String)
	return &c, nil
}

// CreateCasbinPolicy 获取Casbin策略
func (a *CasbinAuther) CreateCasbinPolicy(org string, c *CasbinPolicy) error {
	// log.Println(c.ModelText)
	// 获取角色间的关系
	if rrs, err := new(schema.CasbinGpaRoleRole).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
		// 没有有效的角色关系
	} else if len(*rrs) > 0 {
		for _, v := range *rrs {
			rr := []string{v.ParentName, v.ChildName}

			// 角色前增加应用标识， 标记应用专有角色
			if v.ParentSvc.Valid {
				rr[0] = v.ParentSvc.String + ":" + rr[0]
			}
			if v.ChildSvc.Valid {
				rr[1] = v.ChildSvc.String + ":" + rr[1]
			}
			// 角色前增加Casbin角色专有前缀
			rr[0] = CasbinRolePrefix + rr[0]
			rr[1] = CasbinRolePrefix + rr[1]
			c.Grouping = append(c.Grouping, rr)
		}
	}
	if rps, err := new(schema.CasbinGpaRolePolicy).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
		// 没有有效角色策略关系
	} else if len(*rps) > 0 {
		for _, v := range *rps {
			rp := []string{v.Role, v.Policy}

			// 角色前增加应用标识， 标记应用专有角色
			if v.Svc.Valid {
				rp[0] = v.Svc.String + ":" + rp[0]
			}
			// 角色前增加Casbin角色专有前缀
			rp[0] = CasbinRolePrefix + rp[0]
			// 策略前增加Casbin策略专有前缀
			rp[1] = CasbinPolicyPrefix + rp[1]
			c.Grouping = append(c.Grouping, rp)
		}
	}
	if pss, err := new(schema.CasbinGpaPolicyStatement).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
	} else if len(*pss) > 0 {
		for _, v := range *pss {
			// 策略前增加Casbin策略专有前缀
			sub := CasbinPolicyPrefix + v.Name
			eft := helper.IfString(v.Effect, "allow", "deny")
			if v.Action.Valid {
				actions := strings.Split(v.Action.String, ";")
				for _, action := range actions {
					sa := strings.SplitN(action, ":", 2)
					if len(sa) != 2 {
						break
					}
					svc := sa[0]
					if pas, err := new(schema.CasbinGpaPolicyServiceAction).QueryActionByNameAndSvc(a.Sqlx, sa[1], sa[0]); err != nil {
						if !sqlxc.IsNotFound(err) {
							return err
						}
					} else if len(*pas) > 0 {
						for _, a := range *pas {
							if a.Resource.Valid {
								paths := strings.Split(a.Resource.String, ";")
								for _, path := range paths {
									meth := "*"
									if offset := strings.IndexRune(path, ' '); offset > 0 {
										meth = path[:offset]
										path = path[offset+1:]
									}
									pp := []string{sub, svc, org, path, meth, eft}
									c.Policies = append(c.Policies, pp)
								}
							}
						}
					}
				}
			} else if v.Resource.Valid {
				// 配置资源访问权限， 暂时没有进行开发
			}
		}
	}

	return nil
}

// QueryServiceCode 查询服务
// "fmes:svc-code:" + host + ":" + resource
func (a *CasbinAuther) QueryServiceCode(ctx *gin.Context, user auth.UserInfo, host, path, org string) (string, int, error) {
	if host == "" && path == "" {
		return "", 0, errors.New("无法确认访问的系统服务")
	}
	resource := ""
	if strings.HasPrefix(path, "/api/") {
		rescount := 3
		for i, r := range path {
			if r == '/' {
				if rescount--; rescount == 0 {
					resource = path[:i]
					break
				}
			}
		}
	}
	key := "fmes:svc-code:" + host + ":" + resource

	if svc, b, err := a.Storer.Get(ctx, key); err != nil {
		return "", 0, err // 查询缓存出现异常
	} else if b {
		if strings.HasPrefix(svc, "err:") {
			return "", 0, errors.New(svc[4:]) // 上一次查询，拒绝请求
		}
		offset := strings.IndexRune(svc, '/')
		if offset <= 0 {
			a.Storer.Delete(ctx, key)
			return "", 0, errors.New("系统缓存异常:[" + key + "]" + svc)
		}
		sid, _ := strconv.Atoi(svc[offset+1:])
		return svc[:offset], sid, nil
	}

	// 由于查询是居于全局的， 所以1分钟的缓存是一个合理的范围
	sa := schema.CasbinGpaSvcAud{}
	if err := sa.QueryByAudAndResAndOrg(a.Sqlx, host, resource, ""); err != nil || !sa.SvcCode.Valid {
		// 系统没有配置或者系统为指定有效服务名称
		a.Storer.Set(ctx, key, "err:"+err.Error(), time.Minute) // 1分钟延迟刷新， 拒绝请求也需要缓存
		return "", 0, err
	} else if !sa.SvcCode.Valid {
		a.Storer.Set(ctx, key, "err: no service name", time.Minute)
		return "", 0, errors.New("no service name")
	}
	a.Storer.Set(ctx, key, sa.SvcCode.String+"/"+strconv.Itoa(int(sa.SvcID.Int64)), time.Minute) // 查询结果缓存1分钟
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
		if res == "1" {
			return true, nil
		}
		offset := strings.IndexRune(res, '/')
		if offset <= 0 {
			a.Storer.Delete(ctx, key)
			return false, errors.New("系统缓存异常:[" + key + "]" + res)
		}
		return false, helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: res[:offset], Other: res[offset+1:]})
	}

	var emsg *i18n.Message
	so := schema.CasbinGpaSvcOrg{}
	// 1:启用 0:禁用 2:未激活 3: 删除 4: 欠费 5: 到期
	if err := so.QueryByOrgAndSvc(a.Sqlx, org, sid); err != nil {
		if !sqlxc.IsNotFound(err) {
			return false, err // 系统内部的位置异常
		}
		emsg = &i18n.Message{ID: "WARN-SERVICE-NOFOUND", Other: "访问的服务不存在"}
	} else if so.Expired.Valid && time.Now().After(so.Expired.Time) {
		// 前置授权异常
		emsg = &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"}
	} else if so.Status == schema.StatusEnable {
		// 正常结果返回
		expiration := 2 * time.Minute // 2分钟延迟刷新
		if so.Expired.Valid && so.Expired.Time.Sub(time.Now()) < expiration {
			expiration = so.Expired.Time.Sub(time.Now()) // 修正过期时间
		}
		a.Storer.Set(ctx, key, "1", expiration)
		return true, nil
	} else if so.Status == schema.StatusDisable {
		emsg = &i18n.Message{ID: "WARN-SERVICE-DISABLE", Other: "服务已经被禁用"}
	} else if so.Status == schema.StatusDeleted {
		emsg = &i18n.Message{ID: "WARN-SERVICE-DELETE", Other: "服务已经被删除"}
	} else if so.Status == schema.StatusNoActivate {
		emsg = &i18n.Message{ID: "WARN-SERVICE-NOACTIVATE", Other: "服务未激活"}
	} else if so.Status == schema.StatusExpired {
		emsg = &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"}
	} else {
		emsg = &i18n.Message{ID: "WARN-SERVICE-OTHER", Other: "授权状态异常"}
	}
	a.Storer.Set(ctx, key, emsg.ID+"/"+emsg.Other, time.Minute) // 1分钟延迟刷新， 拒绝请求也需要缓存
	return false, helper.New0Error(ctx, helper.ShowWarn, emsg)
}
