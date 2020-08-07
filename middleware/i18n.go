package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/modules/helper"
)

// I18nMiddleware 国际化
func I18nMiddleware(bundle *i18n.Bundle) gin.HandlerFunc {
	// bundle := i18n.NewBundle(language.Chinese)
	// bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	// bundle.LoadMessageFile("locales/active.zh-CN.toml")
	// bundle.LoadMessageFile("locales/active.en-US.toml")
	return func(c *gin.Context) {
		lang := c.Request.FormValue("lang")
		accept := c.Request.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)
		helper.SetI18n(c, localizer)
		c.Next()
		// 基于i18n句柄可能会跨域生命周期的存在,比如ws协议,所以不应该主动清空
		// helper.SetI18n(c, nil) // 清除

		// localizer.Localize(&i18n.LocalizeConfig{
		// 	DefaultMessage: &i18n.Message{
		// 		ID:    "PersonCatsX",
		// 		Other: "{{.Name}} 有 {{.Count}} 只猫.",
		// 	},
		// 	TemplateData: map[string]interface{}{
		// 		"Name":  "Nick",
		// 		"Count": 2,
		// 	},
		// 	PluralCount: 2,
		// }) // Nick has 2 cats.
	}

}
