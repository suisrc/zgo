package passwd

import (
	"crypto/md5"
	"fmt"
)

// MD5Hash MD5哈希值
func MD5Hash(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// MD5HashString MD5哈希值
func MD5HashString(s string) string {
	return MD5Hash([]byte(s))
}

// VerifyMD5 bcrypt
func VerifyMD5(ent IEntity) (bool, error) {
	epwd := Encrypt([]byte(ent.Left()), []byte(ent.Salt()))
	pwds := MD5Hash(epwd)
	for i := mCost; i > 0; i-- {
		pwds = MD5HashString(pwds)
	}
	return ent.Right() == pwds, nil
}

// GenerateMD5 bcrypt
func GenerateMD5(password string, ptype string) (*GeneratePasswd, error) {
	salt := UUID(32)
	// pwdx := password + "$" + salt // 是否需要拉长密码?进行补位?待定
	epwd := Encrypt([]byte(password), []byte(salt))
	pwds := MD5Hash(epwd)
	for i := mCost; i > 0; i-- {
		pwds = MD5HashString(pwds)
	}
	//verx := "$01$03$"
	return &GeneratePasswd{
		Password:     pwds,
		PasswordType: ptype,
		PasswordSalt: salt,
	}, nil
}
