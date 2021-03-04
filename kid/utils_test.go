package kid_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/kid"
)

func TestGetLocalMac(t *testing.T) {
	macs := kid.GetLocalMac()
	log.Println(macs)

	assert.NotNil(t, nil)
}

func TestGetLocalIP(t *testing.T) {
	ips := kid.GetLocalIP()
	log.Println(ips)

	assert.NotNil(t, nil)
}

func TestGetMustMachineCode(t *testing.T) {
	ips := kid.GetMustMachineCode2()
	log.Println(ips)

	assert.NotNil(t, nil)
}

func TestNewNowCode(t *testing.T) {
	str := kid.NewNowCode(8, false)
	log.Println(str)

	assert.NotNil(t, nil)
}
