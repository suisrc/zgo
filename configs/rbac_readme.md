## 定义
[request_definition]
r = sub, aud, dom, pat, cip, act

[policy_definition]
p = sub, dom, pat, cip, act, eft

## 解释
sub: role, 角色名称, 登陆系统后,一个人只能具有一个角色,如果有多角色,可以使用角色包含来实现
dom: domain, 域, 即当前登陆应用的域名.
pat: path, 路径, 即请求请求访问的资源地址.
cip: client ip, 客户端IP
act: action, 操作(GET)|(POST)|(PUT)|(DELETE)

aud: audience, 基于jwt认证, 通过jwt令牌,获取该令牌签发的接收方, 以确定权限, 需要和 domainAudMatch 一起使用

eft: effect, allow|deny, 允许还是拒绝