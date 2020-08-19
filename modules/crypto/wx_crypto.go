package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"sort"
)

/**
 * 对公众平台发送给公众账号的消息加解密示例代码.
 * 提供基于PKCS7算法的加解
 */
const (
	blockSize = 32
)

// WxPKCS7Encode 获得对明文进行补位填充的字节.
// @param count 需要进行填充补位操作的明文字节个数
// @return 补齐用的字节数组
func WxPKCS7Encode(count int) []byte {
	// 计算需要填充的位数
	amountToPad := blockSize - (count % blockSize)
	if amountToPad == 0 {
		amountToPad = blockSize
	}
	// 获得补位所用的字符
	padChr := WxChr(amountToPad)
	var tmp bytes.Buffer
	for index := 0; index < amountToPad; index++ {
		tmp.WriteRune(padChr)
	}
	return tmp.Bytes()
}

// WxPKCS7Decode 删除解密后明文的补位字符
// @param decrypted 解密后的明文
// @return 删除补位字符后的明文
func WxPKCS7Decode(decrypted []byte) []byte {
	pad := decrypted[len(decrypted)-1]
	if pad < 1 || pad > blockSize {
		pad = 0
	}
	return decrypted[:len(decrypted)-int(pad)]
}

// PKCS7Padding 使用PKCS7进行填充
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padtext...)
}

// PKCS7UnPadding 删除PKCS7填充
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// WxChr 将数字转化成ASCII码对应的字符，用于对明文进行补码
// @param a 需要转化的数字
// @return 转化得到的字符
func WxChr(a int) rune {
	target := rune(a & 0xFF)
	return target
}

// WxByteGroup byte group
type WxByteGroup struct {
	ByteContainer bytes.Buffer
}

// ToBytes to
func (a WxByteGroup) ToBytes() []byte {
	return a.ByteContainer.Bytes()
}

// AddBytes add
func (a WxByteGroup) AddBytes(bytes []byte) {
	a.ByteContainer.Write(bytes)
}

// Size size
func (a WxByteGroup) Size() int {
	return a.ByteContainer.Len()
}

// WxNewCrypto new
func WxNewCrypto(appid, token, encodingAesKey string) *WxCrypto {
	// 必须使用RFC2045标准执行解密
	aesKey, err := Base64DecodeStringMIME(encodingAesKey)
	if err != nil {
		panic(err)
	}
	return &WxCrypto{
		AesKey:         aesKey,
		Token:          token,
		AppID:          appid,
		EncodingAesKey: encodingAesKey,
	}
}

// Number2BytesInNetworkOrder 将一个数字转换成生成4个字节的网络字节序bytes数组
func Number2BytesInNetworkOrder(number int) []byte {
	orderBytes := make([]byte, 4)
	orderBytes[3] = byte(number & 0xFF)
	orderBytes[2] = byte(number >> 8 & 0xFF)
	orderBytes[1] = byte(number >> 16 & 0xFF)
	orderBytes[0] = byte(number >> 24 & 0xFF)
	return orderBytes
}

// BytesNetworkOrder2Number 4个字节的网络字节序bytes数组还原成一个数字
func BytesNetworkOrder2Number(bytesInNetworkOrder []byte) int {
	sourceNumber := 0
	for i := 0; i < 4; i++ {
		sourceNumber <<= 8
		sourceNumber |= int(bytesInNetworkOrder[i]) & 0xFF
	}
	return sourceNumber
}

// WxCrypto wechat 加密
type WxCrypto struct {
	AesKey         []byte
	Token          string
	AppID          string
	EncodingAesKey string
}

// Encrypt 对明文进行加密
// @param plainText 需要加密的明文
// @return 加密后base64编码的字符串
func (a WxCrypto) Encrypt(plainText string) (string, error) {
	randomStr := UUID2(16)

	randomStringBytes := []byte(randomStr)
	plainTextBytes := []byte(plainText)
	bytesOfSizeInNetworkOrder := Number2BytesInNetworkOrder(len(plainTextBytes))
	appIDBytes := []byte(a.AppID)

	byteCollector := WxByteGroup{}

	// randomStr + networkBytesOrder + text + appid
	byteCollector.AddBytes(randomStringBytes)
	byteCollector.AddBytes(plainTextBytes)
	byteCollector.AddBytes(bytesOfSizeInNetworkOrder)
	byteCollector.AddBytes(appIDBytes)

	// ... + pad: 使用自定义的填充方式对明文进行补位填充
	padBytes := WxPKCS7Encode(byteCollector.Size())
	byteCollector.AddBytes(padBytes)

	// 获得最终的字节流, 未加密
	unencrypted := byteCollector.ToBytes()

	//create aes
	cip, err := aes.NewCipher(a.AesKey)
	if err != nil {
		return "", err
	}
	//encrypt string
	cbc := cipher.NewCBCEncrypter(cip, a.AesKey)
	encrypted := make([]byte, len(unencrypted))
	cbc.CryptBlocks(encrypted, unencrypted)

	cipherText := Base64EncodeToString(encrypted)
	return cipherText, nil
}

// Decrypt 对密文进行解密.
// @param cipherText 需要解密的密文
// @param appidOrCorpid 获取解密内容回调，如果为空，会强制判断该内容是否和加密器中的ID相同
// @return 解密得到的明文
func (a WxCrypto) Decrypt(cipherText string, appIDCheck func(string) error) (string, error) {
	cip, err := aes.NewCipher(a.AesKey)
	if err != nil {
		return "", err
	}
	encrypted, err := Base64DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	// CBC mode always works in whole blocks.
	if len(encrypted)%blockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}
	//encrypt string
	cbc := cipher.NewCBCDecrypter(cip, a.AesKey)
	unencrypted := make([]byte, len(encrypted))
	cbc.CryptBlocks(unencrypted, encrypted)

	// 去除补位字符
	content := WxPKCS7Decode(unencrypted)
	// 分离16位随机字符串,网络字节序和AppId
	networkOrder := content[16:24]
	plainTextLen := BytesNetworkOrder2Number(networkOrder)

	appIDBytes := content[20+plainTextLen:]
	appID := string(appIDBytes)

	if appIDCheck != nil {
		if err := appIDCheck(appID); err != nil {
			return "", err
		}
	} else if appID != a.AppID {
		//return "", fmt.Errorf("AppID Error: %s -> %s", appID, a.AppID)
		return "", errors.New("AppID Error")
	}

	plainTextBytes := content[20 : 20+plainTextLen]
	plainText := string(plainTextBytes)
	return plainText, nil
}

// WxGenSHA1 排序,串接arr参数，生成sha1 digest
func WxGenSHA1(arr ...string) string {
	if len(arr) == 0 {
		return ""
	}
	strs := make([]string, len(arr))
	copy(strs, arr)

	var builder bytes.Buffer
	sort.Strings(strs)
	for _, v := range strs {
		builder.WriteString(v)
	}

	// return builder.String()
	return SHA1Hash(builder.Bytes())
}

// WxGenSHA1And 排序,串接arr参数，生成sha1 digest
func WxGenSHA1And(arr ...string) string {
	if len(arr) == 0 {
		return ""
	}
	strs := make([]string, len(arr))
	copy(strs, arr)
	sort.Strings(strs)

	var builder bytes.Buffer
	for _, v := range strs {
		if builder.Len() > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(v)
	}

	// return builder.String()
	return SHA1Hash(builder.Bytes())
}

// WxGenMD5 排序，加密, 用于红包数据签名
func WxGenMD5(datas map[string]string, sign string) string {
	// keys := reflect.ValueOf(datas).MapKeys()
	keys := make([]string, len(datas))
	for k := range datas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder bytes.Buffer
	for _, k := range keys {
		builder.WriteString(k)
		builder.WriteByte('=')
		builder.WriteString(datas[k])
		builder.WriteByte('&')
	}
	builder.WriteString("key=")
	builder.WriteString(sign)

	// return builder.String()
	return SHA1Hash(builder.Bytes())
}
