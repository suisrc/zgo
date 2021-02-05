package oauth2

import (
	"log"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNowTime(t *testing.T) {
	url1 := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=123&agentid=124"

	url2 := url.QueryEscape(url1)
	log.Println(url2)

	assert.NotNil(t, nil)
}

func Test(t *testing.T) {
	token := ""
	openid := WeixinQyOpenid{}
	if err := openid.ConvertToOpenid(token, "15140420103"); err != nil {
		log.Println(err)
	} else {
		log.Println(openid)
	}
	assert.NotNil(t, nil)
}
