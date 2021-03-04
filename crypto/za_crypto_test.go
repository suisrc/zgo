package crypto_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/crypto"
)

func TestAliyunSign(t *testing.T) {
	str := "POS"
	sig := crypto.AliyunSign(str, "AS", "&")

	log.Println(sig)

	assert.NotNil(t, nil)

}
