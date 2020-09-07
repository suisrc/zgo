package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesEncryptBytes aes
func AesEncryptBytes(plainTextBytes, keys []byte) ([]byte, error) {
	// randomStr := UUID2(16)
	// randomStringBytes := []byte(randomStr)
	randomStringBytes := RandomBytes(16)

	var byteCollector bytes.Buffer
	byteCollector.Write(randomStringBytes)
	byteCollector.Write(plainTextBytes)
	unencrypted := PKCS7Padding(byteCollector.Bytes(), 32)
	//create aes
	cip, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}
	//encrypt
	cbc := cipher.NewCBCEncrypter(cip, keys[:cip.BlockSize()])
	encrypted := make([]byte, len(unencrypted))
	cbc.CryptBlocks(encrypted, unencrypted)

	return encrypted, nil
}

// AesEncrypt aes
func AesEncrypt(plainText string, keys []byte) (string, error) {
	plainTextBytes := []byte(plainText)
	encrypted, err := AesEncryptBytes(plainTextBytes, keys)
	if err != nil {
		return "", err
	}

	cipherText := Base64EncodeToString(encrypted)
	return cipherText, nil
}

// AesDecryptBytes aes
func AesDecryptBytes(encrypted, keys []byte) ([]byte, error) {
	cip, err := aes.NewCipher(keys)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(cip, keys[:cip.BlockSize()])
	unencrypted := make([]byte, len(encrypted))
	cbc.CryptBlocks(unencrypted, encrypted)

	// 去除补位字符
	content := PKCS7UnPadding(unencrypted, 32)
	// 删除16位随机数
	return content[16:], nil
}

// AesDecrypt aes
func AesDecrypt(cipherText string, keys []byte) (string, error) {
	encrypted, err := Base64DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	plainTextBytes, err := AesDecryptBytes(encrypted, keys)
	if err != nil {
		return "", err
	}
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

// AesEncryptStr2 aes
func AesEncryptStr2(cipherText, keysStr string) (string, error) {
	keys, err := Base64DecodeString(keysStr)
	if err != nil {
		return "", err
	}
	return AesEncrypt(cipherText, keys)
}
