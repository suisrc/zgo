package crypto_test

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/crypto"
)

func TestT112(t *testing.T) {
	src := "dicksi3你好9D8Dskc_K*&563)()#*(&*D*(ASD&*("
	dst1 := crypto.Base62EncodeToString([]byte(src))
	dst2 := crypto.Base62EncodeToString2([]byte(src))

	log.Println(dst1, "-----", dst2)

	src1, _ := crypto.Base62DecodeString(dst1)
	src2, _ := crypto.Base62DecodeString2(dst2)

	log.Println(string(src1), "-----", string(src2))

	assert.True(t, false)
}

func TestT113(t *testing.T) {
	src := "dicksi3你好9D8Dskc_K*&563)()#*(&*D*(ASD&*("

	log.Println(time.Now())
	for i := 0; i < 10_000_000; i++ {
		crypto.Base62EncodeToString([]byte(src))
	}
	log.Println(time.Now())
	assert.True(t, false)
}

func TestT114(t *testing.T) {
	src := "1234567890"
	dst2 := crypto.Base62EncodeToString3([]byte(src), 'l')

	log.Println(dst2)

	src2, _ := crypto.Base62DecodeString3(dst2, 'l')

	log.Println(string(src2))

	assert.True(t, false)
}
