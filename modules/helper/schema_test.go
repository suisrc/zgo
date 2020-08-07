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
	b := IsResponseError(res)
	assert.True(t, b)

	err := errors.New("hello")
	b = IsResponseError(err)
	assert.False(t, b)

	b = IsResponseError(Err403Forbidden)
	assert.False(t, b)

	b = IsResponseError(ResError(nil, Err403Forbidden))
	assert.True(t, b)
}
