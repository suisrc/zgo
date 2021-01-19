package service

import (
	"bytes"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/modules/crypto"
)

//============================================================================================

// EncryptCaptchaByAccount 加密验证码
func EncryptCaptchaByAccount(c *gin.Context, account int, secret, captcha string, expire time.Duration) (string, error) {
	expireAt := time.Now().Add(expire).Unix()                      // 过期时间
	var buffer bytes.Buffer                                        // byte buffer
	buffer.Write(crypto.Number2BytesInNetworkOrder(int(expireAt))) // 占用4个字节
	buffer.Write(crypto.Number2BytesInNetworkOrder(account))       // 占用4个字节
	buffer.Write([]byte(captcha))                                  // 验证码
	// 加密KEY的内容
	keys, err := crypto.Base64DecodeString(secret)
	if err != nil {
		return "", err
	}
	// 对验证码进行加密, 验证码后端不存储, 加密成校验码
	checkCodeBytes, err := crypto.AesEncryptBytes(buffer.Bytes(), keys)
	if err != nil {
		return "", err // 加密出现问题
	}
	// 给出登陆账户信息,以用来进行解密
	resultCodeBytes := append(checkCodeBytes, crypto.Number2BytesInNetworkOrder(account)...)
	return crypto.Base64EncodeToString(resultCodeBytes), nil
}

// DecryptCaptchaByAccount 解密验证码
func DecryptCaptchaByAccount(c *gin.Context, code string) ( /*account*/ int, func( /*secret*/ string) (string, time.Duration, error), error) {
	resultCodeBytes, err := crypto.Base64DecodeString(code)
	if err != nil {
		return 0, nil, err
	}
	resultCodeLen := len(resultCodeBytes)
	if resultCodeLen < 12 {
		// 不可预知异常, 往往来自恶意攻击
		return 0, nil, errors.New("code is error")
	}
	account := crypto.BytesNetworkOrder2Number(resultCodeBytes[resultCodeLen-4:])
	return account, func(secret string) (string, time.Duration, error) {
		keys, err := crypto.Base64DecodeString(secret)
		if err != nil {
			return "", 0, err
		}
		checkCodeBytes, err := crypto.AesDecryptBytes(resultCodeBytes[:resultCodeLen-4], keys)
		if err != nil {
			return "", 0, err
		}
		account2 := crypto.BytesNetworkOrder2Number(checkCodeBytes[4:8])
		if account != account2 {
			// 不可预知异常, 往往来自恶意攻击
			return "", 0, errors.New("account id error")
		}
		expireAt := crypto.BytesNetworkOrder2Number(checkCodeBytes[:4])
		expire := time.Unix(int64(expireAt), 0).Sub(time.Now())
		captcha := string(checkCodeBytes[8:])
		return captcha, expire, nil
	}, nil
}

//============================================================================================

// EncryptAccountWithUser 加密验证码
func EncryptAccountWithUser(c *gin.Context, account, user int, secret string) (string, error) {
	var buffer bytes.Buffer                                  // byte buffer
	buffer.Write(crypto.Number2BytesInNetworkOrder(account)) // 占用4个字节
	buffer.Write(crypto.Number2BytesInNetworkOrder(user))    // 占用4个字节
	buffer.Write(crypto.RandomBytes(4))                      // 掩码4个字节
	// 加密KEY的内容
	keys := []byte(crypto.FixRandomAes32(secret))
	// 加密
	resultCodeBytes, err := crypto.AesEncryptBytes(buffer.Bytes(), keys)
	if err != nil {
		return "", err // 加密出现问题
	}
	return crypto.Base64EncodeToStringURL(resultCodeBytes), nil
}

// DecryptAccountWithUser 解密验证码
func DecryptAccountWithUser(c *gin.Context, data string, secret string) (account, user int, err error) {
	dataBytes, err := crypto.Base64DecodeStringURL(data)
	if err != nil {
		return 0, 0, err
	}
	resultCodeLen := len(dataBytes)
	if resultCodeLen < 12 {
		return 0, 0, errors.New("data is error") // 不可预知异常, 往往来自恶意攻击
	}
	keys := []byte(crypto.FixRandomAes32(secret))
	resultCodeBytes, err := crypto.AesDecryptBytes(dataBytes, keys)
	if err != nil {
		return 0, 0, err
	}
	account = crypto.BytesNetworkOrder2Number(resultCodeBytes[:4])
	user = crypto.BytesNetworkOrder2Number(resultCodeBytes[4:8])
	return
}
