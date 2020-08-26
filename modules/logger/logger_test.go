package logger

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack3(t *testing.T) {
	info := Stack3(1, 2, false)
	log.Println(fmt.Sprintf("%s", info))

	assert.True(t, true)
}

func TestWhereis(t *testing.T) {
	info := Whereis(true)
	log.Println(fmt.Sprintf("%s", info))

	assert.True(t, true)
}

func TestErrorWW(t *testing.T) {
	ShowShortName = true
	info := ErrorWW(errors.New("error"))
	log.Println(fmt.Sprintf("%s", info))

	assert.True(t, true)
}
