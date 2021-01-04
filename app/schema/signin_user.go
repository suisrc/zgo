package schema

import "github.com/suisrc/zgo/modules/auth"

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	TokenID   string
	UserID    string
	Name      string
	Role      string
	XidxID    string
	AccountID string
	T3rdID    string
	ClientID  string
	Domain    string
	Issuer    string
	Audience  string
	OrgCode   string
	OrgRole   string
	OrgDomain string
	OrgAdmin  string
}

// GetTokenID token
func (u *SigninUser) GetTokenID() string {
	return u.TokenID
}

// GetUserID user
func (u *SigninUser) GetUserID() string {
	return u.UserID
}

// GetUserName name
func (u *SigninUser) GetUserName() string {
	return u.Name
}

// GetUserRole role
func (u *SigninUser) GetUserRole() string {
	return u.Role
}

// SetUserRole role
func (u *SigninUser) SetUserRole(nrole string) string {
	orole := u.Role
	u.Role = nrole
	return orole
}

// GetXidxID user index
func (u *SigninUser) GetXidxID() string {
	return u.XidxID
}

// GetAccountID token
func (u *SigninUser) GetAccountID() string {
	return u.AccountID
}

// GetT3rdID 3rd id
func (u *SigninUser) GetT3rdID() string {
	return u.T3rdID
}

// GetClientID client id
func (u *SigninUser) GetClientID() string {
	return u.ClientID
}

// GetDomain domain
func (u *SigninUser) GetDomain() string {
	return u.Domain
}

// GetIssuer issuer
func (u *SigninUser) GetIssuer() string {
	return u.Issuer
}

// GetAudience audience
func (u *SigninUser) GetAudience() string {
	return u.Audience
}

// GetOrgCode org code
func (u *SigninUser) GetOrgCode() string {
	return u.OrgCode
}

// GetOrgRole org role
func (u *SigninUser) GetOrgRole() string {
	return u.OrgRole
}

// GetOrgDomain org domain
func (u *SigninUser) GetOrgDomain() string {
	return u.OrgDomain
}

// GetOrgAdmin org admin
func (u *SigninUser) GetOrgAdmin() string {
	return u.OrgAdmin
}
