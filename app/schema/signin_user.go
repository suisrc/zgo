package schema

import (
	"strings"

	"github.com/suisrc/zgo/modules/auth"
)

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	TokenID   string
	AccountID string
	UserIdxID string
	UserID    string
	UserName  string
	UserRoles string
	OrgCode   string
	OrgAdmin  string
	Domain    string
	Issuer    string
	Audience  string
}

// GetTokenID xxx
func (u *SigninUser) GetTokenID() string {
	return u.TokenID
}

// GetAccountID xxx
func (u *SigninUser) GetAccountID() string {
	return u.AccountID
}

// GetUserIdxID xxx
func (u *SigninUser) GetUserIdxID() string {
	return u.UserIdxID
}

// GetUserID xxx
func (u *SigninUser) GetUserID() string {
	return u.UserID
}

// GetUserName xxx
func (u *SigninUser) GetUserName() string {
	return u.UserName
}

// GetUserRoles xxx
func (u *SigninUser) GetUserRoles() []string {
	return strings.Split(u.UserRoles, ";")
}

// GetOrgCode xxx
func (u *SigninUser) GetOrgCode() string {
	return u.OrgCode
}

// GetOrgAdmin xxx
func (u *SigninUser) GetOrgAdmin() string {
	return u.OrgAdmin
}

// GetDomain xxx
func (u *SigninUser) GetDomain() string {
	return u.Domain
}

// GetIssuer xxx
func (u *SigninUser) GetIssuer() string {
	return u.Issuer
}

// GetAudience xxx
func (u *SigninUser) GetAudience() string {
	return u.Audience
}

// GetUserSvcRoles xxx
func (u *SigninUser) GetUserSvcRoles(svc string) []string {
	roles := []string{}
	for _, role := range u.GetUserRoles() {
		if strings.HasPrefix(role, svc) {
			roles = append(roles, role)
		}
	}
	return roles
}

// SetUserRoles xxx
func (u *SigninUser) SetUserRoles(roles []string) {
	u.UserRoles = strings.Join(roles, ";")
}
