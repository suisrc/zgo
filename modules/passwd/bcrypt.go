package passwd

import (
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// VerifyBcrypt bcrypt
func VerifyBcrypt(ent IEntity) (bool, error) {
	hashpass := ent.Salt() + ent.Right()
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(ent.Left()))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GenerateBcrypt bcrypt
func GenerateBcrypt(password string, ptype string) (*GeneratePasswd, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bCost)
	if err != nil {
		return nil, err
	}
	pwdstr := string(pwd)
	// $ver$cost$[salte:22]hashpass
	offset := strings.LastIndex(pwdstr, "$") + 22
	return &GeneratePasswd{
		Password:     pwdstr[offset:],
		PasswordType: ptype,
		PasswordSalt: pwdstr[:offset],
	}, nil
}

// VerifyBcrypt2 bcrypt
func VerifyBcrypt2(ent IEntity) (bool, error) {
	salt, err := base64.StdEncoding.DecodeString(ent.Salt())
	if err != nil {
		return false, nil
	}
	hashpass := reverse(string(salt)) + ent.Right()
	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(ent.Left()))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GenerateBcrypt2 bcrypt
func GenerateBcrypt2(password string, ptype string) (*GeneratePasswd, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return nil, err
	}
	pwdstr := string(pwd)
	// $ver$cost$[salte:22]hashpass
	offset := strings.LastIndex(pwdstr, "$") + 22
	salt := base64.StdEncoding.EncodeToString([]byte(reverse(pwdstr[:offset])))
	return &GeneratePasswd{
		Password:     pwdstr[offset:],
		PasswordType: ptype,
		PasswordSalt: salt,
	}, nil
}

// VerifyBcrypt3 bcrypt
func VerifyBcrypt3(ent IEntity) (bool, error) {
	salx, err := base64.StdEncoding.DecodeString(ent.Salt())
	if err != nil {
		return false, nil
	}
	sbyt := Decrypt(salx, []byte(ent.Right()))
	salt := string(sbyt)
	hashpass := salt + ent.Right()
	err = bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(ent.Left()))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GenerateBcrypt3 bcrypt
func GenerateBcrypt3(password string, ptype string) (*GeneratePasswd, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 9)
	if err != nil {
		return nil, err
	}
	pwdstr := string(pwd)
	// $ver$cost$[salte:22]hashpass
	offset := strings.LastIndex(pwdstr, "$") + 22
	salt := pwdstr[:offset]
	rpwd := pwdstr[offset:]
	sbyt := Encrypt([]byte(salt), []byte(rpwd))
	salx := base64.StdEncoding.EncodeToString(sbyt)
	return &GeneratePasswd{
		Password:     rpwd,
		PasswordType: ptype,
		PasswordSalt: salx,
	}, nil
}
