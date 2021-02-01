package schema

import (
	"strings"

	"github.com/suisrc/zgo/modules/auth"
)

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	TokenID   string
	Account   string
	UserID    string
	UserName  string
	UserRoles string
	OrgCode   string
	OrgAdmin  string
	OrgUsrID  string
	OrgAppID  string
	Domain    string
	Issuer    string
	Audience  string
	CustomID  string
}

// GetTokenID xxx
func (u *SigninUser) GetTokenID() string {
	return u.TokenID
}

// GetAccount xxx
func (u *SigninUser) GetAccount() string {
	return u.Account
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
	if u.UserRoles == "" {
		return nil
	}
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

// GetOrgUsrID xxx
func (u *SigninUser) GetOrgUsrID() string {
	return u.OrgUsrID
}

// GetOrgAppID xxx
func (u *SigninUser) GetOrgAppID() string {
	return u.OrgAppID
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
	if roles == nil {
		u.UserRoles = ""
	} else {
		u.UserRoles = strings.Join(roles, ";")
	}
}
