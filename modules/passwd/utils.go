package passwd

import "github.com/NebulousLabs/fastrand"

func reverse(str string) string {
	bytes := []rune(str)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	str = string(bytes)
	return str
}

// Encrypt 简单的加密算法, 破坏其本身字符串的特征, 当用于对同一个MD5值进行处理
func Encrypt(buffers []byte, randoms []byte) []byte {
	lenBuf := len(buffers)
	lenEnd := lenBuf / 2
	for index := 0; index < lenEnd; index += 2 {
		offset := lenBuf - 1 - index
		buffers[index] ^= ^buffers[offset] // Go语言取反方式和C语言不同，Go语言不支持~符号
		buffers[offset] ^= ^buffers[index]
		buffers[index] ^= ^buffers[offset]
	}
	lenRan := len(randoms)
	for index, offset := 0, 0; index < lenBuf; index++ {
		if offset >= lenRan {
			offset = 0
		}
		buf := randoms[offset]
		buffers[index] = (^buffers[index] ^ ^buf)
		offset++
	}
	return buffers
}

// Decrypt 简单的解密算法
func Decrypt(buffers []byte, randoms []byte) []byte {
	lenBuf := len(buffers)
	lenRan := len(randoms)
	for index, offset := 0, 0; index < lenBuf; index++ {
		if offset >= lenRan {
			offset = 0
		}
		buf := randoms[offset]
		buffers[index] = ^(^buffers[index] ^ buf)
		offset++
	}
	lenEnd := lenBuf / 2
	for index := 0; index < lenEnd; index += 2 {
		offset := lenBuf - 1 - index
		buffers[index] ^= ^buffers[offset]
		buffers[offset] ^= ^buffers[index]
		buffers[index] ^= ^buffers[offset]
	}
	return buffers
}

// UUID uuid
func UUID(length int64) string {
	ele := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	elen := len(ele)
	uuid := ""
	var i int64
	for i = 0; i < length; i++ {
		uuid += ele[fastrand.Intn(elen)]
	}
	return uuid
}
