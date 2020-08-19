package crypto

import "strings"

// Reverse 字符串倒序
func Reverse(str string) string {
	bytes := []rune(str)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	str = string(bytes)
	return str
}

// FixPreStrLen2 获取固定长度内容，不够前段补齐0, 超长, 保留原始数据
func FixPreStrLen2(pstr string, plen int) string {
	slen := len(pstr)
	if slen > plen || slen == plen {
		return pstr
	}
	var sbir strings.Builder
	count := plen - len(pstr)
	for i := 0; i < count; i++ {
		sbir.WriteByte('0')
	}
	sbir.WriteString(pstr)
	return sbir.String()
}

// FixPreStrLen 获取固定长度内容，不够前段补齐0
func FixPreStrLen(pstr string, plen int) string {
	slen := len(pstr)
	if slen > plen {
		return pstr[:plen]
	}
	if slen == plen {
		return pstr
	}
	var sbir strings.Builder
	count := plen - len(pstr)
	for i := 0; i < count; i++ {
		sbir.WriteByte('0')
	}
	sbir.WriteString(pstr)
	return sbir.String()
}

// FixSufStrLen 获取固定长度内容，不够前段补齐0
func FixSufStrLen(pstr string, plen int) string {
	slen := len(pstr)
	if slen > plen {
		return pstr[:plen]
	}
	if slen == plen {
		return pstr
	}
	var sbir strings.Builder
	sbir.WriteString(pstr)
	count := plen - len(pstr)
	for i := 0; i < count; i++ {
		sbir.WriteByte('0')
	}
	return sbir.String()
}
