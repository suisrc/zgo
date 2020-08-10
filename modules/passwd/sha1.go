package passwd

import (
	"crypto/sha1"
	"fmt"
)

// SHA1Hash SHA1哈希值
func SHA1Hash(b []byte) string {
	h := sha1.New()
	_, _ = h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1HashString SHA1哈希值
func SHA1HashString(s string) string {
	return SHA1Hash([]byte(s))
}

// VerifySHA1 bcrypt
func VerifySHA1(ent IEntity) (bool, error) {
	epwd := Encrypt([]byte(ent.Left()), []byte(ent.Salt()))
	pwds := SHA1Hash(epwd)
	for i := mCost; i > 0; i-- {
		pwds = SHA1HashString(pwds)
	}
	return ent.Right() == pwds, nil
}

// GenerateSHA1 bcrypt
func GenerateSHA1(password string, ptype string) (*GeneratePasswd, error) {
	salt := UUID(32)
	epwd := Encrypt([]byte(password), []byte(salt))
	pwds := SHA1Hash(epwd)
	for i := mCost; i > 0; i-- {
		pwds = SHA1HashString(pwds)
	}
	return &GeneratePasswd{
		Password:     pwds,
		PasswordType: ptype,
		PasswordSalt: salt,
	}, nil
}
