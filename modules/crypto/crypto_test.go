package crypto

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestT1(t *testing.T) {
	log.Println(UUID2(12))
	assert.True(t, true)
}

func TestCrpt(t *testing.T) {
	string1 := UUID(16)
	log.Println("加密:" + string1)
	buffers := []byte(string1)
	log.Println(buffers)
	// randoms := []byte(uuid(32))
	randoms := []byte("987654")

	encfers := MaskEncrypt(buffers, randoms)
	log.Println(encfers)
	log.Println(Base64EncodeToString(encfers))

	decfers := MaskDecrypt(encfers, randoms)

	string2 := string(decfers)
	log.Println("解密:" + string2)

	assert.True(t, string1 == string2)
}

func TestBase64(t *testing.T) {
	s0 := "lBXYSlGJuQcFPiS4KCfLGxQjmcHJRrJuoIfrKC2NPwt"
	s1 := Base64EncodeToString([]byte(s0))
	log.Println(s1)
	s2, _ := Base64DecodeString(s1)
	log.Println(string(s2))

	s3, _ := Base64DecodeString(s0 + "=")
	log.Println(len(s3))
	log.Println(s3)

	s4, _ := Base64DecodeString(s0 + "=")
	log.Println(s4)

	assert.True(t, true)
}

func TestAes32(t *testing.T) {
	s0 := RandomAes32()
	log.Println(s0)

	s2, _ := Base64DecodeString(s0)
	log.Println(s2)

	assert.True(t, true)
}

func TestWxCrypto(t *testing.T) {
	aesKey0 := RandomAes32()
	log.Println(aesKey0)

	aesKey1, _ := Base64DecodeString(aesKey0)

	s0 := "kdixkdiskdiDiskc"

	s1, err1 := AesEncrypt(s0, aesKey1)
	assert.Nil(t, err1)

	log.Println(s1)

	s2, err2 := AesDecrypt(s1, aesKey1)
	assert.Nil(t, err2)

	log.Println(s2)

	assert.Equal(t, s0, s2)
}
