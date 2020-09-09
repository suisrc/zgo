package oauth2

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/schema"
)

var _ Handler = &WeixinQy{}

// WeixinQy 微信企业号
type WeixinQy struct {
}

// Handle handle
func (a *WeixinQy) Handle(*gin.Context, *schema.SigninOfOAuth2, *schema.SigninGpaOAuth2Platfrm, *schema.SigninGpaAccount) error {
	return errors.New("未实现")
}
