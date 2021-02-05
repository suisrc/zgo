package injector

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/api"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/app"
	"golang.org/x/text/language"
)

// InjectorSet 注入Injector
var InjectorSet = wire.NewSet(
	NewBundle, // 国际化注册

	wire.Struct(new(service.I18n), "*"), // i18n数据库依赖
	service.InitI18nLoader,              // i18n数据库依赖
)

// InjectorEndSet 注入Injector
var InjectorEndSet = wire.NewSet(
	// app.NewSwagger,                  // swagger
	app.NewHealthz,                  // 健康检查
	wire.Struct(new(Injector), "*"), // 注册器
)

//======================================
// 注入控制器
//======================================

// Injector 注入器(用于初始化完成之后的引用)
type Injector struct {
	Engine    *gin.Engine    // gin引擎
	Endpoints *api.Endpoints // api接口

	Bundle     *i18n.Bundle       // 国际化
	I18nLoader service.I18nLoader // i18n 数据库加载器

	// Swagger app.Swagger
	Healthz app.Healthz
}

//======================================
// END
//======================================

// NewBundle 国际化
func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.Chinese)
	//bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("locales/active.zh-CN.toml")
	bundle.LoadMessageFile("locales/active.en-US.toml")
	return bundle
}
