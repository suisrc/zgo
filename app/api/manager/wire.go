package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// EndpointSet wire注入声明
var EndpointSet = wire.NewSet(
	wire.Struct(new(User), "*"),
	wire.Struct(new(Account), "*"),
	wire.Struct(new(Role), "*"),
	wire.Struct(new(Menu), "*"),
	wire.Struct(new(Gateway), "*"),
)

// Wire 注入控制器
type Wire struct {
	User    *User
	Account *Account
	Role    *Role
	Menu    *Menu
	Gateway *Gateway
}

// Register 主路由必须包含UAC内容
func (a *Wire) Register(r gin.IRouter) {
	m := r.Group("manager")

	a.User.Register(m)
	a.Account.Register(m)
	a.Role.Register(m)
	a.Menu.Register(m)
	a.Gateway.Register(m)
}
