package crypto

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWxCrypto2(t *testing.T) {
	// wc := WxNewCrypto("123456", "IDKxiddis98", RandomAes32())
	wc := WxNewCrypto2("123456", "IDKxiddis98", "lBXYSlGJuQcFPiS4KCfLGxQjmcHJRrJuoIfrKC2NPwt")
	log.Println(wc.EncodingAesKey)

	text := "你好, golang, {}, IDixudDLSOCKSIcskDI, DNIs /slo ////*sd*(<xml?>"
	etext, err := wc.Encrypt(text)
	assert.Nil(t, err)
	log.Println(etext)

	utext, err := wc.Decrypt(etext)
	assert.Nil(t, err)
	log.Println(utext)

}

func TestBytesNetworkOrder2Number(t *testing.T) {
	n := Number2BytesInNetworkOrder(14)
	log.Println(n)

	b := BytesNetworkOrder2Number(n)
	log.Println(b)

	assert.True(t, true)
}
