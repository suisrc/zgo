package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strings"
)

// aliyun 加密

// AliyunSignature ...
/**
 * @see https://github.com/alibabacloud-go/openapi-util/blob/c929a595b5b8321f44ceee2afcd84ea744f56df6/service/service.go#L475
 * @param signedParams params which need to be signed
 * @param method http method e.g. GET
 * @param secret AccessKeySecret
 * @return the signature
 */
func AliyunSignature(signedParams map[string]string, method string, secret string) string {
	stringToSign := aliyunBuildStringToSign(signedParams, method)
	signature := AliyunSign(stringToSign, secret, "&")
	return signature
}

func aliyunBuildStringToSign(signedParam map[string]string, method string) (stringToSign string) {
	signParams := make(map[string]string)
	for key, value := range signedParam {
		signParams[key] = value
	}
	delete(signParams, "Signature")

	stringToSign = aliyunGetFormedMap(signParams)
	stringToSign = strings.Replace(stringToSign, "+", "%20", -1)
	stringToSign = strings.Replace(stringToSign, "*", "%2A", -1)
	stringToSign = strings.Replace(stringToSign, "%7E", "~", -1)
	stringToSign = url.QueryEscape(stringToSign)
	stringToSign = method + "&%2F&" + stringToSign
	return
}

func aliyunGetFormedMap(source map[string]string) (urlEncoded string) {
	urlEncoder := url.Values{}
	for key, value := range source {
		urlEncoder.Add(key, value)
	}
	urlEncoded = urlEncoder.Encode()
	return
}

// AliyunSign ...
func AliyunSign(stringToSign, accessKeySecret, secretSuffix string) string {
	// log.Println(stringToSign)
	secret := accessKeySecret + secretSuffix
	signedBytes := aliyunShaHmac1(stringToSign, secret)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func aliyunShaHmac1(source, secret string) []byte {
	key := []byte(secret)
	hmac := hmac.New(sha1.New, key)
	hmac.Write([]byte(source))
	return hmac.Sum(nil)
}
