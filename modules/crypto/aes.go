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

// AesEncryptStr aes
func AesEncryptStr(plainText string) (string, string, error) {
	keys := RandomBytes(32)
	cipherText, err := AesEncrypt(plainText, keys)
	if err != nil {
		return "", "", err
	}
	return cipherText, Base64EncodeToString(keys), nil
}

// AesDecryptStr aes
func AesDecryptStr(cipherText, keysStr string) (string, error) {
	keys, err := Base64DecodeString(keysStr)
	if err != nil {
		return "", err
	}
	return AesDecrypt(cipherText, keys)
}
