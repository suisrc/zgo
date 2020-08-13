package middleware

import (
	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
)

// I18nMiddleware 国际化
func I18nMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	// bundle := i18n.NewBundle(language.Chinese)
	// bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// bundle.LoadMessageFile("locales/active.zh-CN.toml")
	// bundle.LoadMessageFile("locales/active.en-US.toml")
	return i18n.Serve(bundle)
}
