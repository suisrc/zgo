package helper

// IfString 判断
func IfString(condition bool, ifture, ifalse string) string {
	if condition {
		return ifture
	}
	return ifalse
}

// IfInt 判断
func IfInt(condition bool, ifture, ifalse int) int {
	if condition {
		return ifture
	}
	return ifalse
}

// IfBool 判断
func IfBool(condition bool, ifture, ifalse bool) bool {
	if condition {
		return ifture
	}
	return ifalse
}

// IfObject 判断
func IfObject(condition bool, ifture, ifalse interface{}) interface{} {
	if condition {
		return ifture
	}
	return ifalse
}

// IfFunc 判断
func IfFunc(condition bool, ifture, ifalse func() interface{}) interface{} {
	if condition {
		return ifture()
	}
	return ifalse()
}

// ReverseStr 反正字符串
func ReverseStr(s string) string {
	if len(s) <= 1 {
		return s
	}
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// SplitStrCR 截取字符串
func SplitStrCR(s string, x rune, c int) string {
	for i, r := range s {
		if r == x {
			if c--; c == 0 {
				return s[:i]
			}
		}
	}
	return s
}
