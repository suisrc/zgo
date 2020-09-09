package jwt

import (
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/crypto"
)

// NewUserInfo 获取用户信息
func NewUserInfo(user auth.UserInfo) *UserClaims {
	claims := UserClaims{}

	tokenID := user.GetTokenID()
	if tokenID == "" {
		if ati, err := strconv.Atoi(user.GetAccountID()); err != nil {
			tokenID = NewRandomID(user.GetAccountID())
		} else {
			tokenID = NewRandomID(crypto.EncodeBaseX32(int64(ati)))
		}
	}

	claims.Id = tokenID
	claims.Subject = user.GetUserID()
	claims.Name = user.GetUserName()
	claims.Role = user.GetRoleID()

	claims.Issuer = user.GetIssuer()
	claims.Audience = user.GetAudience()

	claims.AccountID = user.GetAccountID()

	return &claims
}

var _ auth.UserInfo = &UserClaims{}

// UserClaims 用户信息声明
type UserClaims struct {
	jwt.StandardClaims
	Name       string      `json:"nam,omitempty"` // 用户名
	Role       string      `json:"rol,omitempty"` // 角色ID, role id
	AccountID  string      `json:"ati,omitempty"` // 登陆ID, 本身不具备任何意义,只是标记登陆方式
	Properties interface{} `json:"pps,omitempty"` // 用户的额外属性
}

// GetUserName name
func (u *UserClaims) GetUserName() string {
	return u.Name
}

// GetUserID user
func (u *UserClaims) GetUserID() string {
	return u.Subject
}

// GetRoleID role
func (u *UserClaims) GetRoleID() string {
	// if u.Role == "" {
	// 	if u.Subject == "1" {
	// 		return "admin" // 作为默认系统用户
	// 	}
	// 	return "default" // 当用户没有权限的时候,启动默认权限
	// }
	return u.Role
}

// SetRoleID role
func (u *UserClaims) SetRoleID(nrole string) string {
	orole := u.Role
	u.Role = nrole
	return orole
}

// GetTokenID token
func (u *UserClaims) GetTokenID() string {
	return u.Id
}

// GetAccountID token
func (u *UserClaims) GetAccountID() string {
	return u.AccountID
}

// GetIssuer issuer
func (u *UserClaims) GetIssuer() string {
	return u.Issuer
}

// GetAudience audience
func (u *UserClaims) GetAudience() string {
	return u.Audience
}

// GetProps props
func (u *UserClaims) GetProps() (interface{}, bool) {
	return nil, false
}
