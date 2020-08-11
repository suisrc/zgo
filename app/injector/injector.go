package injector

import (
	"github.com/BurntSushi/toml"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/api"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/middlewire"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/auth/jwt/store/buntdb"
	casbinzgo "github.com/suisrc/zgo/modules/casbin"
	casbinmem "github.com/suisrc/zgo/modules/casbin/watcher/mem"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/logger"
	"golang.org/x/text/language"
)

// InjectorSet 注入Injector
var InjectorSet = wire.NewSet(
	NewBundle, // 国际化注册
	NewAuther, // Auther注册

	//casbinjson.CasbinAdapterSet, // Casbin依赖
	service.CasbinAdapterSet, // Casbin依赖
	wire.Bind(new(casbinzgo.PolicyVer), new(service.CasbinAdapter)), // Casbin版本
	casbinmem.NewCasbinWatcher,                                      // Casbin观察者

	wire.Struct(new(service.I18n), "*"), // i18n数据库依赖
	service.InitI18nLoader,              // i18n数据库依赖
)

// InjectorEndSet 注入Injector
var InjectorEndSet = wire.NewSet(
	middlewire.NewHealthz,           // 健康检查
	middlewire.NewSwagger,           // swagger
	wire.Struct(new(Injector), "*"), // 注册器
)

//======================================
// 注入控制器
//======================================

// Injector 注入器(用于初始化完成之后的引用)
type Injector struct {
	Engine    *gin.Engine    // gin引擎
	Endpoints *api.Endpoints // api接口

	Bundle     *i18n.Bundle           // 国际化
	Enforcer   *casbin.SyncedEnforcer // 权限认证
	Auther     auth.Auther            // 令牌控制
	Watcher    persist.Watcher        // casbin观察者
	I18nLoader service.I18nLoader     // i18n 数据库加载器

	Swagger middlewire.Swagger
	Healthz middlewire.Healthz
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

// NewAuther of auth.Auther
// 授权认证使用的auther内容
func NewAuther() auth.Auther {
	store, err := buntdb.NewStore(":memory:") // 使用内存缓存
	if err != nil {
		panic(err)
	}
	secret := config.C.JWTAuth.SigningSecret
	if secret == "" {
		secret = auth.UUID(128)
		logger.Infof(nil, "jwt secret: %s", secret)
	}
	auther := jwt.New(store,
		jwt.SetSigningSecret(secret), // 注册令牌签名密钥
	)

	return auther
}
