package kid

import (
	"strings"
	"sync"
	"time"

	"github.com/NebulousLabs/fastrand"
	"github.com/suisrc/zgo/crypto"
)

var shareCnt uint64 = fastrand.Uint64n(1 << 63)
var shareLck sync.Mutex

// NewSequenceCode ... 排序性质
func NewSequenceCode(size int, x62 bool) string {
	if size <= 0 {
		return ""
	}
	shareLck.Lock()
	shareCnt++
	idx := int64(shareCnt) // 得到新计数器
	shareLck.Unlock()
	return NewIdxCode(size, idx, x62)
}

// NewIdxCode ... 编码
func NewIdxCode(size int, idx int64, x62 bool) string {
	if size <= 0 {
		return ""
	}
	code := ""
	if x62 {
		code = crypto.EncodeBaseX62(idx)
	} else {
		code = crypto.EncodeBaseX32(idx)
	}
	if csize := len(code); csize > size {
		return code[csize-size:]
	} else if csize == size {
		return code
	} else {
		builder := strings.Builder{}
		for i := 0; i < size-csize; i++ {
			builder.WriteRune('0') // 时间补码
		}
		builder.WriteString(code)
		return builder.String()
	}
}

// NewRandomCode ...
func NewRandomCode(size int, x64 bool) string {
	if x64 {
		return crypto.UUID2(size)
	}
	return crypto.UUID(size)
}

// NewNowCode 获取当前时间编码
func NewNowCode(size int, x64 bool) string {
	if size <= 0 {
		return ""
	}
	code := ""
	if x64 {
		code = crypto.EncodeBaseX64(time.Now().Unix())
	} else {
		code = crypto.EncodeBaseX32(time.Now().Unix())
	}
	if csize := len(code); csize > size {
		return code[csize-size:]
	} else if csize == size {
		return code
	} else {
		builder := strings.Builder{}
		for i := 0; i < size-csize; i++ {
			builder.WriteRune('0') // 时间补码
		}
		builder.WriteString(code)
		return builder.String()
	}
}
