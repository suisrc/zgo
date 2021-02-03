package module

import "time"

// CasbinObject subject
type CasbinObject struct {
	Svc    string
	Host   string
	Path   string
	Method string
	Client string
}

// CasbinSubject subject
type CasbinSubject struct {
	//UsrID    int
	//AccID    int
	Acc    string
	Org    string
	Usr    string
	OrgUsr string
	OrgApp string
	Iss    string
	Aud    string
	Role   string
	Scope  string
}

var (
	// CasbinPolicyModel casbin使用的对比模型
	CasbinPolicyModel = `[request_definition]
r = sub, obj

[policy_definition]
p = sub, svc, org, path, meth, eft, c8n

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = `
	// CasbinDefaultMatcher casbin使用的对比模型
	CasbinDefaultMatcher = `g(r.sub.Role, p.sub) && r.obj.Svc==p.svc && (p.meth=="" || methodMatch(r.obj.Method, p.meth)) && (p.path=="" || keyMatch(r.obj.Path, p.path)) && (p.c8n=="" || customMatch(p.c8n, r.sub, r.obj))`
)

var (
	// CasbinCachedExpireAt 缓存定时器刷新时间
	CasbinCachedExpireAt = 4 * time.Minute
	// CasbinEnforcerCheckAt 引擎检测版本时间
	CasbinEnforcerCheckAt = 2 * time.Minute
	// CasbinEnforcerExpireAt 引擎标记过期时间
	CasbinEnforcerExpireAt = 8 * time.Minute
	// CasbinServiceCodeExpireAt 服务缓存过期时间
	CasbinServiceCodeExpireAt = 2 * time.Minute
	// CasbinServiceTenantExpireAt 租户缓存过期时间
	CasbinServiceTenantExpireAt = 2 * time.Minute
)

var (
	// CasbinSvcRoleKey 角色配置
	CasbinSvcRoleKey = "X-Request-Svc-[SVC-NAME]-Role"
	// CasbinSysRoleKey 系统平台角色
	CasbinSysRoleKey = "X-Request-Sys-Role"
	// CasbinSvcPublic 公共服务
	CasbinSvcPublic = "pub-"
	// CasbinRolePrefix 角色
	CasbinRolePrefix = "r:"
	// CasbinUserPrefix 用户
	CasbinUserPrefix = "u:"
	// CasbinPolicyPrefix 策略
	CasbinPolicyPrefix = "p:"
	// // CasbinNoSign 未登陆
	// CasbinNoSign = "nosign"
	// // CasbinNoRole 无角色
	// CasbinNoRole = "norole"
	// // CasbinNoUser 无用户
	// CasbinNoUser = "nouser"
)
