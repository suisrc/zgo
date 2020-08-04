package casbin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	b, err := DomainMatchFunc("", "a.cn", "a.cn")
	assert.Nil(t, err)
	assert.True(t, b.(bool))

	b, err = DomainMatchFunc("", "a.a.cn", "*.a.cn")
	assert.Nil(t, err)
	assert.True(t, b.(bool))

	b, err = DomainMatchFunc(".a.cn", "a.a.cn", "jwt")
	assert.Nil(t, err)
	assert.True(t, b.(bool))

	b, err = DomainMatchFunc(".a.cn", "a.a.cn", "")
	assert.Nil(t, err)
	assert.True(t, b.(bool))
}
