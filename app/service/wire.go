package service

import (
	"github.com/google/wire"
	"github.com/suisrc/zgo/app/model/gpaf"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/app/oauth2"
	"github.com/suisrc/zgo/modules/passwd"
)

// ServiceSet wire注入服务
var ServiceSet = wire.NewSet(
	gpaf.NewGPA,                             // 数据库链接
	wire.Struct(new(passwd.Validator), "*"), // 密码验证
	oauth2.NewSelector,                      // OAuth2注册

	module.NewAuther,   // Auther注册
	module.NewStorer,   // Storer注册
	module.NewEventBus, // 事件总线
	wire.Struct(new(module.AuthOpts), "GPA", "Storer"),               // Auther依赖
	wire.Struct(new(module.CasbinAuther), "GPA", "Storer", "Auther"), // Casbin依赖

	wire.Struct(new(MobileSender), "*"), // 手机发送器
	wire.Struct(new(EmailSender), "*"),  // 邮件发送器
	wire.Struct(new(ThreeSender), "*"),  // 第三方平台发送器

	wire.Struct(new(Signin), "*"), // 登陆服务
	wire.Struct(new(User), "*"),   // 用户服务
)

//======================================
// 分割线
//======================================

// ResultRef 返回值暂存器
type ResultRef struct {
	D interface{}
}
