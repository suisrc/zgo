package service

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
	// UsrID    int
	// AccID    int
	Org    string
	Usr    string
	OrgUsr string
	OrgApp string
	Iss    string
	Aud    string
}

var (
	// CasbinPolicyModel casbin使用的对比模型
	CasbinPolicyModel = `[request_definition]
r = sub, obj, role

[policy_definition]
p = sub, svc, org, path, meth, eft

[role_definition]
g = _, _
g2 = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = `
	// CasbinDefaultMatcher casbin使用的对比模型
	CasbinDefaultMatcher = `g(r.role, p.sub) && (p.meth=="" || methodMatch(r.obj.Method, p.meth)) && (p.path=="" || keyMatch(r.obj.Path, p.path))`
	// `(g(r.role, p.sub) || keyMatch(p.sub, "u:*") && g2(r.sub.Usr, p.sub)) && (p.path=="" || keyMatch(r.obj.Path, p.path))`
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
