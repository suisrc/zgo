package helper

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// 定义上下文中的键
const (
	ResI18nKey = Prefix + "/res-i18n"
)

// FormatMessage fm
func FormatMessage(c *gin.Context, lc *i18n.LocalizeConfig) string {
	return MustI18n(c).MustLocalize(lc)
}

// FormatCode fc
func FormatCode(c *gin.Context, message *i18n.Message, args map[string]interface{}) string {
	if localizer, ok := GetI18n(c); ok {
		return localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: message,
			TemplateData:   args,
		})
	}
	// 加载i18n后,不会进入该分支
	if args == nil {
		return message.Other
	}
	text := message.Other
	for key, val := range args {
		text = strings.ReplaceAll(text, "{{."+key+"}}", toString(val))
	}
	return text
}

// MustI18n 用户
func MustI18n(c *gin.Context) *i18n.Localizer {
	localizer, ok := GetI18n(c)
	if !ok {
		panic(errors.New("context no has i18n localizer"))
	}
	return localizer
}

// GetI18n 用户
func GetI18n(c *gin.Context) (*i18n.Localizer, bool) {
	if v, ok := c.Get(ResI18nKey); ok {
		if l, b := v.(*i18n.Localizer); b {
			return l, true
		}
	}
	return nil, false
}

// SetI18n 用户
func SetI18n(c *gin.Context, l *i18n.Localizer) {
	c.Set(ResI18nKey, l)
}

// ToString to string
func toString(a interface{}) string {
	switch a.(type) {
	case int:
		return strconv.Itoa(a.(int))
	case bool:
		return strconv.FormatBool(a.(bool))
	case string:
		return a.(string)
	default:
		return "conver error"
	}
}
