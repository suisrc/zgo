package crypto_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/crypto"
)

func TestT111(t *testing.T) {
	log.Println(crypto.UUID2(12))
	assert.True(t, true)
}

func TestCrpt(t *testing.T) {
	string1 := crypto.UUID(16)
	log.Println("加密:" + string1)
	buffers := []byte(string1)
	log.Println(buffers)
	// randoms := []byte(uuid(32))
	randoms := []byte("987654")

	encfers := crypto.MaskEncrypt(buffers, randoms)
	log.Println(encfers)
	log.Println(crypto.Base64EncodeToString(encfers))

	decfers := crypto.MaskDecrypt(encfers, randoms)

	string2 := string(decfers)
	log.Println("解密:" + string2)

	assert.True(t, string1 == string2)
}

func TestBase64(t *testing.T) {
	s0 := "lBXYSlGJuQcFPiS4KCfLGxQjmcHJRrJuoIfrKC2NPwt"
	s1 := crypto.Base64EncodeToString([]byte(s0))
	log.Println(s1)
	s2, _ := crypto.Base64DecodeString(s1)
	log.Println(string(s2))

	s3, _ := crypto.Base64DecodeString(s0 + "=")
	log.Println(len(s3))
	log.Println(s3)

	s4, _ := crypto.Base64DecodeString(s0 + "=")
	log.Println(s4)

	assert.True(t, true)
}

func TestAes32(t *testing.T) {
	s0 := crypto.RandomAes32()
	log.Println(s0)

	s2, _ := crypto.Base64DecodeString(s0)
	log.Println(s2)

	assert.True(t, true)
}

func TestCrypto(t *testing.T) {
	aesKey0 := crypto.RandomAes32()
	log.Println(aesKey0)

	aesKey1, _ := crypto.Base64DecodeString(aesKey0)

	s0 := "kdixkdiskdiDiskc"

	s1, err1 := crypto.AesEncrypt(s0, aesKey1)
	assert.Nil(t, err1)

	log.Println(s1)

	s2, err2 := crypto.AesDecrypt(s1, aesKey1)
	assert.Nil(t, err2)

	log.Println(s2)

	assert.Equal(t, s0, s2)
}

func TestCrypto2(t *testing.T) {
	s0 := "kdixkdiskdiDiskc"

	s1, key, _ := crypto.AesEncryptStr(s0)

	log.Println(key)
	log.Println(s1)

	s2, _ := crypto.AesDecryptStr(s1, key)

	log.Println(s2)

	assert.Equal(t, s0, s2)
}
