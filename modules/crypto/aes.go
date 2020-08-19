package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesEncrypt aes
func AesEncrypt(plainText string, keys []byte) (string, error) {
	// randomStr := UUID2(16)
	// randomStringBytes := []byte(randomStr)
	randomStringBytes := RandomBytes(16)
	plainTextBytes := []byte(plainText)

	var byteCollector bytes.Buffer
	byteCollector.Write(randomStringBytes)
	byteCollector.Write(plainTextBytes)
	unencrypted := PKCS7Padding(byteCollector.Bytes(), 32)
	//create aes
	cip, err := aes.NewCipher(keys)
	if err != nil {
		return "", err
	}
	//encrypt string
	cbc := cipher.NewCBCEncrypter(cip, keys[:cip.BlockSize()])
	encrypted := make([]byte, len(unencrypted))
	cbc.CryptBlocks(encrypted, unencrypted)

	cipherText := Base64EncodeToString(encrypted)
	return cipherText, nil
}

// AesDecrypt aes
func AesDecrypt(cipherText string, keys []byte) (string, error) {
	encrypted, err := Base64DecodeString(cipherText)

	cip, err := aes.NewCipher(keys)
	if err != nil {
		return "", err
	}
	cbc := cipher.NewCBCDecrypter(cip, keys[:cip.BlockSize()])
	unencrypted := make([]byte, len(encrypted))
	cbc.CryptBlocks(unencrypted, encrypted)

	// 去除补位字符
	content := PKCS7UnPadding(unencrypted, 32)

	plainTextBytes := content[16:]
	return string(plainTextBytes), nil
}
