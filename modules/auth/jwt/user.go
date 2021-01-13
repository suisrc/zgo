package jwt

import (
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/crypto"
)

// NewTokenID new ID
func NewTokenID(_ati string) string {
	var builder strings.Builder
	if ati, err := strconv.Atoi(_ati); err != nil {
		builder.WriteString(_ati)
	} else {
		builder.WriteString(crypto.EncodeBaseX32(int64(ati)))
	}
	builder.WriteRune('t')
	builder.WriteString(crypto.UUID(20))                         // 20位
	builder.WriteString(crypto.EncodeBaseX32(time.Now().Unix())) // 7位
	return builder.String()
}

// NewRefreshToken new refresh token
func NewRefreshToken(_ati string) string {
	var builder strings.Builder
	if ati, err := strconv.Atoi(_ati); err != nil {
		builder.WriteString(_ati)
	} else {
		builder.WriteString(crypto.EncodeBaseX32(int64(ati)))
	}
	builder.WriteRune('r')
	builder.WriteString(crypto.UUID(20))                         // 20位
	builder.WriteString(crypto.EncodeBaseX32(time.Now().Unix())) // 7位
	builder.WriteString(crypto.UUID(12))                         // 12位
	return builder.String()
}

// NewUserInfo 获取用户信息
func NewUserInfo(user auth.UserInfo) *UserClaims {
	claims := UserClaims{}

	tokenID := user.GetTokenID()
	if tokenID == "" {
		tokenID = NewTokenID(user.GetAccountID())
	}
	claims.Id = tokenID
	claims.AccountID = user.GetAccountID()
	claims.UserIdxID = user.GetUserIdxID()

	claims.Subject = user.GetUserID()
	claims.UserName = user.GetUserName()
	claims.UserRoles = user.GetUserRoles()

	claims.OrgCode = user.GetOrgCode()
	claims.OrgAdmin = user.GetOrgAdmin()

	claims.Domain = user.GetDomain()
	claims.Issuer = user.GetIssuer()
	claims.Audience = user.GetAudience()

	return &claims
}

var _ auth.UserInfo = &UserClaims{}

// UserClaims 用户信息声明
type UserClaims struct {
	jwt.StandardClaims

	// TokenID -> Id
	AccountID string `json:"ati,omitempty"` // 登陆ID, 本身不具备任何意义,只是标记登陆方式, 使用token反向加密
	UserIdxID string `json:"idx,omitempty"` // 用户的一种扩展ID, 为用户索引ID, 使用token反向加密

	// UserID -> Subject -> sub, GetOrgCode为空，提供用户平台ID， 否则提供用户租户ID
	UserName  string   `json:"nam,omitempty"` // 用户名
	UserRoles []string `json:"ros,omitempty"` // 角色ID, 该角色是平台角色， 也可以理解为平台给机构的角色

	OrgCode  string `json:"ogc,omitempty"` // 组织code
	OrgAdmin string `json:"oga,omitempty"` // admin'为用户管理员， GetOrgCode为空，提供

	Domain string `json:"dom,omitempty"` // 业务域，主要用户当前用户跨应用的业务关联，暂时不使用
	// Issuer -> Issuer
	// Audience -> Audience

	// _TmpRoles []string // `json:"-"`
}

// GetTokenID xxx
func (u *UserClaims) GetTokenID() string {
	return u.Id
}

// GetAccountID xxx
func (u *UserClaims) GetAccountID() string {
	return u.AccountID
}

// GetUserIdxID xxx
func (u *UserClaims) GetUserIdxID() string {
	return u.UserIdxID
}

// GetUserID xxx
func (u *UserClaims) GetUserID() string {
	return u.Subject
}

// GetUserName xxx
func (u *UserClaims) GetUserName() string {
	return u.UserName
}

// GetUserRoles xxx
func (u *UserClaims) GetUserRoles() []string {
	return u.UserRoles
}

// GetOrgCode xxx
func (u *UserClaims) GetOrgCode() string {
	return u.OrgCode
}

// GetOrgAdmin xxx
func (u *UserClaims) GetOrgAdmin() string {
	return u.OrgAdmin
}

// GetDomain xxx
func (u *UserClaims) GetDomain() string {
	return u.Domain
}

// GetIssuer xxx
func (u *UserClaims) GetIssuer() string {
	return u.Issuer
}

// GetAudience xxx
func (u *UserClaims) GetAudience() string {
	return u.Audience
}

// GetUserSvcRoles xxx
func (u *UserClaims) GetUserSvcRoles(svc string) []string {
	roles := []string{}
	for _, role := range u.GetUserRoles() {
		if strings.HasPrefix(role, svc) {
			roles = append(roles, role)
		}
	}
	return roles
}
