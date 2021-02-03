package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
)

// Weixin 接口
type Weixin struct {
	gpa.GPA
}

// Register 注册路由
func (a *Weixin) Register(r gin.IRouter) {
	// weixin := r.Group("weixin")
	// r.GET("wx", a.wx1)
	// r.POST("wx", a.wx1)
}

func (a *Weixin) wx1(c *gin.Context) {

	query := &crypto.WxEncryptSignature{}
	if err := helper.ParseQuery(c, query); err != nil {
		c.String(200, "success")
		return // 解析发生异常， 但是需要通知微信服务器接收到请求
	}
	body := &crypto.WxEncryptMessage{}
	if err := helper.ParseJSON(c, body); err != nil {
		c.String(200, "success")
		return // 解析发生异常， 但是需要通知微信服务器接收到请求
	}

	logger.Infof(c, body.Encrypt)

	wc := crypto.WxNewCrypto2("123456", "IDKxiddis98", "lBXYSlGJuQcFPiS4KCfLGxQjmcHJRrJuoIfrKC2NPwt")
	content, err := wc.Decrypt(body.Encrypt)
	if err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
		c.String(200, "success")
		return // 解析发生异常， 但是需要通知微信服务器接收到请求
	}
	logger.Infof(c, content)
	result, err := wc.Encrypt(content)
	if err != nil {
		logger.Errorf(c, logger.ErrorWW(err))
		c.String(200, "success")
		return // 解析发生异常， 但是需要通知微信服务器接收到请求
	}

	em := crypto.WxEncryptMessage{
		Encrypt:   result,
		Nonce:     crypto.UUID(16),
		TimeStamp: strconv.Itoa(int(time.Now().UnixNano() / 1e6)),
	}
	// 生成签名
	em.MsgSignature = crypto.WxGenSHA1(wc.Token, em.TimeStamp, em.Nonce, em.Encrypt)

	c.JSON(200, em)
}
