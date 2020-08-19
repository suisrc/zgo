package crypto

// MaskEncrypt 简单的加密算法, 破坏其本身字符串的特征, 当用于对同一个MD5值进行处理
func MaskEncrypt(buffers []byte, randoms []byte) []byte {
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

// MaskDecrypt 简单的解密算法
func MaskDecrypt(buffers []byte, randoms []byte) []byte {
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
