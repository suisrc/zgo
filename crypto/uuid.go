package crypto

import (
	"bytes"
	"errors"
	"strings"

	"github.com/NebulousLabs/fastrand"
)

var (
	code10 = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

	code36 = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

	code62 = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

	code64 = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'-', '_'}
)

// CODE code
func CODE(length int) string {
	elen := len(code10)
	var uuid strings.Builder
	for i := 0; i < length; i++ {
		uuid.WriteByte(code36[fastrand.Intn(elen)])
	}
	return uuid.String()
}

// UUID uuid
func UUID(length int) string {
	elen := len(code36)
	var uuid strings.Builder
	for i := 0; i < length; i++ {
		uuid.WriteByte(code36[fastrand.Intn(elen)])
	}
	return uuid.String()
}

// UUID2 uuid
func UUID2(length int) string {
	elen := len(code62)
	var uuid strings.Builder
	for i := 0; i < length; i++ {
		uuid.WriteByte(code62[fastrand.Intn(elen)])
	}
	return uuid.String()
}

// EncodeBaseX62 ... z被定义位转义字符
func EncodeBaseX62(code int64) string {
	if code == 0 {
		return "0"
	}
	var sbir strings.Builder
	value := uint64(code) // go 没有无符号左移,需要使用无符号字符处理
	for value != 0 {
		current := value & 0x3F
		if current >= 61 {
			sbir.WriteByte(code64[current-61])
			sbir.WriteRune('z')
		} else {
			sbir.WriteByte(code64[current])
		}
		value >>= 6
	}
	return Reverse(sbir.String())
}

// EncodeBaseX64 64位编码
func EncodeBaseX64(code int64) string {
	if code == 0 {
		return "0"
	}
	var sbir strings.Builder
	value := uint64(code) // go 没有无符号左移,需要使用无符号字符处理
	for value != 0 {
		current := value & 0x3F
		sbir.WriteByte(code64[current])
		value >>= 6
	}
	return Reverse(sbir.String())
}

// DecodeBaseX64 64解码
func DecodeBaseX64(code string) (int64, error) {
	if strings.TrimSpace(code) == "" || code == "0" {
		return 0, nil
	}
	var value int64
	for _, ch := range []byte(code) {
		if '0' <= ch && ch <= '9' {
			value = (value << 6) + (int64(ch) - int64('0'))
		} else if 'A' <= ch && ch <= 'Z' {
			value = (value << 6) + (int64(ch) - int64('A') + 10)
		} else if 'a' <= ch && ch <= 'z' {
			value = (value << 6) + (int64(ch) - int64('a') + 36)
		} else if '-' == ch {
			value = (value << 6) + 62
		} else if '_' == ch {
			value = (value << 6) + 63
		} else {
			return 0, errors.New(string(ch) + ":code is must in [0-9A-Za-z-_]")
		}
	}
	return value, nil
}

// EncodeBaseX32 64位编码
func EncodeBaseX32(code int64) string {
	if code == 0 {
		return "0"
	}
	var sbir strings.Builder
	value := uint64(code) // go 没有无符号左移,需要使用无符号字符处理
	for value != 0 {
		current := value & 0x1F
		sbir.WriteByte(code36[current])
		value >>= 5
	}
	return Reverse(sbir.String())
}

// DecodeBaseX32 64解码
func DecodeBaseX32(code string) (int64, error) {
	if strings.TrimSpace(code) == "" || code == "0" {
		return 0, nil
	}
	var value int64
	for _, ch := range []byte(code) {
		if '0' <= ch && ch <= '9' {
			value = (value << 6) + (int64(ch) - int64('0'))
		} else if 'a' <= ch && ch <= 'v' {
			value = (value << 6) + (int64(ch) - int64('a') + 10)
		} else {
			return 0, errors.New(string(ch) + ":code is must in [0-9a-v]")
		}
	}
	return value, nil
}

// RandomBytes random
func RandomBytes(length int) []byte {
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteByte(byte(fastrand.Intn(255)))
	}
	return buffer.Bytes()
}

// RandomAes32 random aes 32
func RandomAes32() string {
	return Base64EncodeToString(RandomBytes(32))
}

// FixRandomAes32 random aes 32
func FixRandomAes32(str string) string {
	return FixPreStrLen(str, 32)
}
