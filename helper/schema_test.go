package helper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsResponseError(t *testing.T) {
	res := &Success{
		Success: true,
	}
	b := FixResponseError(nil, res)
	assert.True(t, b)

	err := errors.New("hello")
	b = FixResponseError(nil, err)
	assert.False(t, b)

	b = FixResponseError(nil, Err403Forbidden)
	assert.False(t, b)

	b = FixResponseError(nil, ResError(nil, Err403Forbidden))
	assert.True(t, b)
}
