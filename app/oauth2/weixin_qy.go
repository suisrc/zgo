package oauth2

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/store"
)

var _ Handler = &WeixinQy{}

// WeixinQy 微信企业号
type WeixinQy struct {
	gpa.GPA              // 数据库操作
	Storer  store.Storer // 缓存器操作
}

// Handle handle
func (a *WeixinQy) Handle(*gin.Context, *schema.SigninOfOAuth2, *schema.SigninGpaOAuth2Platfrm, *schema.SigninGpaAccount) error {
	return errors.New("未实现")
}
