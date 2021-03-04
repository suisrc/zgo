package app

// 对gin进行初始化
import (
	"github.com/suisrc/zgo/config"
	"github.com/suisrc/zgo/middleware"

	"github.com/gin-gonic/gin"
)

// UseEngine 修正Engine内容
type UseEngine func(*gin.Engine)

// InitGinEngine engine
func InitGinEngine(opt UseEngine) *gin.Engine {
	gin.SetMode(config.C.RunMode)
	//gin.SetMode(gin.DebugMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	opt(app)
	//app.Use(gin.Logger())
	//app.Use(middleware.LoggerMiddleware())
	//app.Use(gin.Recovery())
	//app.Use(middleware.RecoveryMiddleware())

	return app
}
