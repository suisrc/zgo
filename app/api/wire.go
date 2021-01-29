package api

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/api/manager"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/middleware"
	"github.com/suisrc/zgo/middlewire"
	"github.com/suisrc/zgo/modules/config"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// EndpointSet wire注入声明
var EndpointSet = wire.NewSet(
	NewUseEngine,                   // 增加引擎中间件
	service.ServiceSet,             // 系统提供的服务列表
	wire.Struct(new(Options), "*"), // 初始化接口参数
	InitEndpoints,                  // 初始化接口方法
	manager.EndpointSet,            // 管理页面接口

	// 接口注册
	wire.Struct(new(Demo), "*"),
	wire.Struct(new(Auth), "*"),
	wire.Struct(new(Signin), "*"),
	wire.Struct(new(User), "*"),
	wire.Struct(new(System), "*"),

	// 界面管理
	wire.Struct(new(manager.Wire), "*"),
)

//=====================================
// Endpoint
//=====================================

// Options options
type Options struct {
	Engine *gin.Engine       // 服务器
	Router middlewire.Router // 根路由

	// 接口注入
	Demo   *Demo
	Auth   *Auth
	Signin *Signin
	User   *User
	System *System

	// 权限管理
	CasbinAuther *service.CasbinAuther

	// 管理界面
	ManagerWire *manager.Wire
}

// Endpoints result
type Endpoints struct {
}

// InitEndpoints init
func InitEndpoints(o *Options) *Endpoints {
	// 在根路由注册通用授权接口, (没有ContextPath限定,一般是给nginx使用)
	// 在nginx注册认证接口时候,请放行zgo服务器其他接口,防止重复认证
	// 注意，改接口为内容接口，为提供国际化语言支持
	o.Auth.Register(o.Engine)

	// ContextPath路由
	r := o.Router
	// 国际化，根路由国际化
	// r.Use(middleware.I18nMiddleware(o.Bundle))

	//router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	//router.GET("/index", func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "index.tmpl", gin.H{
	//		"title": "Main website",
	//	})
	//})

	// 服务器授权控制器
	// 增加权限认证
	uac := o.CasbinAuther.UserAuthCasbinMiddleware(
		middleware.AllowPathPrefixSkipper(
			// sign 登陆接口需要排除
			// 注意[/api/sign,都会被排除]
			middleware.JoinPath(config.C.HTTP.ContextPath, "sign"),
			// pub => public 为系统公共信息
			// 注意[/api/pub开头的,都会被排除]
			middleware.JoinPath(config.C.HTTP.ContextPath, "pub"),
			// /user/current
			middleware.JoinPath(config.C.HTTP.ContextPath, "user", "current"),
		),
	)
	r.Use(uac)

	// 注册登陆接口
	o.Signin.Register(r)
	o.User.Register(r)

	o.System.Register(r)
	o.Demo.Register(o.Engine)

	o.ManagerWire.Register(r)
	return &Endpoints{}
}

// NewUseEngine 绑定中间件
func NewUseEngine(bundle *i18n.Bundle) middlewire.UseEngine {
	return func(app *gin.Engine) {
		app.Use(gin.Logger())
		//app.Use(middleware.LoggerMiddleware())
		//app.Use(gin.Recovery())
		app.Use(middleware.RecoveryMiddleware())
		// 国际化, 全部国际化
		app.Use(middleware.I18nMiddleware(bundle))
	}
}
