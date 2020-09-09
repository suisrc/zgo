package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/store"
)

// Selector 选择器
type Selector map[string]Handler

// NewSelector 全局缓存
func NewSelector(GPA gpa.GPA, Storer store.Storer) (Selector, error) {
	selector := make(map[string]Handler)

	selector["WX"] = &WeixinQm{
		GPA:    GPA,
		Storer: Storer,
	}
	selector["WXQ"] = &WeixinQy{
		GPA:    GPA,
		Storer: Storer,
	}

	return selector, nil
}

// Handler 认证接口
type Handler interface {
	// Handle 处理OAuth2认证
	// account 用户账户
	// domain 请求域, 如果不存在,直接指定"", 其作用是在多应用授权时候,准确定位子应用
	// client 请求端, 定位子应用
	Handle(*gin.Context, *schema.SigninOfOAuth2, *schema.SigninGpaOAuth2Platfrm, *schema.SigninGpaAccount) error
}
