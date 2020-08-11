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
	zgocasbin "github.com/suisrc/zgo/modules/casbin"
	casbinmem "github.com/suisrc/zgo/modules/casbin/watcher/mem"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/logger"
	"golang.org/x/text/language"
)

// InjectorSet 注入Injector
var InjectorSet = wire.NewSet(
	middlewire.NewHealthz, // 健康检查
	NewBundle,             // 国际化注册
	NewAuther,             // Auther注册

	wire.Bind(new(zgocasbin.PolicyVer), new(service.CasbinAdapter)), // Casbin版本
	service.CasbinAdapterSet,   // Casbin依赖
	casbinmem.NewCasbinWatcher, // Casbin观察者
	//casbinjson.CasbinAdapterSet, // Casbin依赖
)

// InjectorEndSet 注入Injector
var InjectorEndSet = wire.NewSet(
	middlewire.NewSwagger,           // swagger
	wire.Struct(new(Injector), "*"), // 注册器
)

//======================================
// 注入控制器
//======================================

// Injector 注入器(用于初始化完成之后的引用)
type Injector struct {
	Engine    *gin.Engine
	Endpoints *api.Endpoints
	Swagger   middlewire.Swagger
	Healthz   middlewire.Healthz

	Bundle   *i18n.Bundle           // 国际化
	Enforcer *casbin.SyncedEnforcer // 权限认证
	Auther   auth.Auther            // 令牌控制
	Watcher  persist.Watcher        // casbin adapter
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
