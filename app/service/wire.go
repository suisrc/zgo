package service

import (
	"github.com/google/wire"
	"github.com/suisrc/zgo/app/model/entc"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/oauth2"
	"github.com/suisrc/zgo/modules/passwd"
)

// ServiceSet wire注入服务
var ServiceSet = wire.NewSet(
	NewAuther,       // Auther注册
	NewStorer,       // Storer注册
	entc.NewClient,  // 数据库连接注册
	sqlxc.NewClient, // 数据库连接注册

	oauth2.NewSelector, // OAuth2注册

	wire.Struct(new(gpa.GPA), "*"),                            // 数据库服务
	wire.Struct(new(passwd.Validator), "*"),                   // 密码验证
	wire.Struct(new(AuthOpts), "GPA", "Storer"),               // Auther依赖
	wire.Struct(new(CasbinAuther), "GPA", "Storer", "Auther"), // Casbin依赖

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
