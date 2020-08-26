package passwd

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/modules/crypto"
)

func TestBcrypt(t *testing.T) {
	pwd, err := GenerateBcrypt("123456", "BCR")
	assert.Nil(t, err)
	byt, err := json.Marshal(pwd)
	log.Println(string(byt))
	pwx := &PasswdX{
		GeneratePasswd: *pwd,
		PasswordX:      "123456",
	}
	res, err := VerifyBcrypt(pwx)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestBcrypt2(t *testing.T) {
	pwd, err := GenerateBcrypt2("123456", "BCR")
	assert.Nil(t, err)
	byt, err := json.Marshal(pwd)
	log.Println(string(byt))
	pwx := &PasswdX{
		GeneratePasswd: *pwd,
		PasswordX:      "123456",
	}
	res, err := VerifyBcrypt2(pwx)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestBcrypt3(t *testing.T) {
	pwd, err := GenerateBcrypt3("123456", "BCR")
	assert.Nil(t, err)
	byt, err := json.Marshal(pwd)
	log.Println(string(byt))
	pwx := &PasswdX{
		GeneratePasswd: *pwd,
		PasswordX:      "123456",
	}
	res, err := VerifyBcrypt3(pwx)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestMD5(t *testing.T) {
	pwd, err := GenerateMD5("123456", "MD5")
	assert.Nil(t, err)
	byt, err := json.Marshal(pwd)
	log.Println(string(byt))
	pwx := &PasswdX{
		GeneratePasswd: *pwd,
		PasswordX:      "123456",
	}
	res, err := VerifyMD5(pwx)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestSHA1(t *testing.T) {
	pwd, err := GenerateSHA1("123456", "SHA1")
	assert.Nil(t, err)
	byt, err := json.Marshal(pwd)
	log.Println(string(byt))
	pwx := &PasswdX{
		GeneratePasswd: *pwd,
		PasswordX:      "123456",
	}
	res, err := VerifySHA1(pwx)
	assert.Nil(t, err)
	assert.True(t, res)
}

type PasswdX struct {
	GeneratePasswd
	PasswordX string
}

func (a *PasswdX) Target() string {
	return a.PasswordX
}

// 测试速度
func TestMD5Speed(t *testing.T) {

	for i := 10000; i > 0; i-- {
		pwo := crypto.UUID(16)
		pwd, _ := GenerateMD5(pwo, "MD5")
		pwx := &PasswdX{
			GeneratePasswd: *pwd,
			PasswordX:      pwo,
		}
		VerifyMD5(pwx)
	}

	assert.True(t, true)
}

// 测试速度
func TestSHA1Speed(t *testing.T) {

	for i := 10000; i > 0; i-- {
		pwo := crypto.UUID(16)
		pwd, _ := GenerateSHA1(pwo, "SHA1")
		pwx := &PasswdX{
			GeneratePasswd: *pwd,
			PasswordX:      pwo,
		}
		VerifySHA1(pwx)
	}

	assert.True(t, true)
}

// 测试速度
func TestBCR0Speed(t *testing.T) {

	for i := 10; i > 0; i-- {
		pwo := crypto.UUID(16)
		pwd, _ := GenerateBcrypt(pwo, "BCR")
		pwx := &PasswdX{
			GeneratePasswd: *pwd,
			PasswordX:      pwo,
		}
		VerifyBcrypt(pwx)
	}

	assert.True(t, true)
}
