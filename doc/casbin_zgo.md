# 当前系统casbin说明

## 使用的model
```conf
[request_definition]
r = sub, aud, dom, pat, cip, act

[policy_definition]
p = sub, dom, pat, cip, act, eft

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && actionMatch(r.act, p.act) && domainAudMatch(r.dom, p.dom, r.aud) && keyMatch(r.pat, p.pat) && ipMatch(r.cip, p.cip)
```

## 中间件

middleware.UserAuthCasbinMiddleware

## 参数说明
sub: role,        角色名称, 登陆系统后,一个人只能具有一个角色,如果有多角色,可以使用角色包含来实现  
dom: domain,      域, 即当前登陆应用的域名  
pat: path,        路径, 即请求请求访问的资源地址  
cip: client ip,   客户端IP  
act: action,      操作(GET)|(POST)|(PUT)|(DELETE)  
  
aud: audience,    基于jwt认证, 通过jwt令牌,获取该令牌签发的接收方, 以确定权限, 需要和 domainAudMatch 一起使用  

eft: effect,      allow|deny, 允许还是拒绝  

## 函数说明(附加函数)

domainMatch:     验证域名, ""和"*"代表全匹配  
domainAudMatch   验证域名, ""和"*"代表全匹配, "jwt"需要验证[request_definition]中的dom和aud是否相同  
actionMatch      验证操作, ""和"*"代表全匹配, 可以使用regexMatch(r.act, p.act)代替,但是这里使用更简单的Contains方法进行处理  

## 特殊说明
当[policy_definition]中存在sub = [nosignin, norole]时候,会自动激活无角色认证配置, 即config.casbin.nosignin和config.casbin.norole, 而忽略之前的配置  
该配置具有强制性,所以,如果系统拒绝该配置,请在policy中禁止进行[nosignin, norole]资源的配置系统系统中希望有公共资源访问配置,
目前使用"/api/pub*"和"/api/sign*"请求路径的请求都会被强制回来认证,不会强制获取登陆人员信息,用户需要通过auth.Auther获取  