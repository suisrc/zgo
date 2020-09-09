package oauth2

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
)

var _ Handler = (*WeixinQm)(nil)

// WeixinQm 微信公众号
type WeixinQm struct {
}

// Handle handle
func (a *WeixinQm) Handle(*gin.Context, *schema.SigninOfOAuth2, *schema.SigninGpaOAuth2Platfrm, *schema.SigninGpaAccount) error {
	return errors.New("未实现")
}
