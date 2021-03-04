package kid_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/kid"
)

func TestNewRoleKID(t *testing.T) {
	for i := 0; i < 1000; i++ {
		kid := kid.NewRoleKID()
		log.Println(kid)
	}
	assert.NotNil(t, nil)
}
