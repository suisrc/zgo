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
	builder.WriteRune('_')
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
	builder.WriteRune('_')
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
		tokenID = NewTokenID(claims.AccountID)
	}
	claims.Id = tokenID
	claims.Subject = user.GetUserID()

	claims.Name = user.GetUserName()
	claims.Role = user.GetUserRole()
	claims.XidxID = user.GetXidxID()
	claims.AccountID = user.GetAccountID()
	claims.T3rdID = user.GetT3rdID()
	claims.ClientID = user.GetClientID()

	claims.Domain = user.GetDomain()
	claims.Issuer = user.GetIssuer()
	claims.Audience = user.GetAudience()

	claims.OrgCode = user.GetOrgCode()
	claims.OrgRole = user.GetOrgRole()
	claims.OrgDomain = user.GetOrgDomain()
	claims.OrgAdmin = user.GetOrgAdmin()

	return &claims
}

var _ auth.UserInfo = &UserClaims{}

// UserClaims 用户信息声明
type UserClaims struct {
	jwt.StandardClaims
	Name      string `json:"nam,omitempty"` // 用户名
	Role      string `json:"rol,omitempty"` // 角色ID, 该角色是平台角色， 也可以理解为平台给机构的角色
	XidxID    string `json:"xid,omitempty"` // 用户的一种扩展ID, 为用户索引ID
	AccountID string `json:"ati,omitempty"` // 登陆ID, 本身不具备任何意义,只是标记登陆方式
	T3rdID    string `json:"aki,omitempty"` // 子应用用户ID
	ClientID  string `json:"cki,omitempty"` // 子应用应用ID
	Domain    string `json:"dom,omitempty"` // 业务域，主要用户当前用户跨应用的业务关联，暂时不使用
	OrgCode   string `json:"ogc,omitempty"` // 组织code
	OrgRole   string `json:"ogr,omitempty"` // 组织role
	OrgDomain string `json:"ogd,omitempty"` // 组织领域，在OrdCode上进行扩展，细化组织下部门的概念
	OrgAdmin  string `json:"oga,omitempty"` // 是否为管理员，平台没有管理员的概念, 管理员只服务员组织
}

// GetTokenID token
func (u *UserClaims) GetTokenID() string {
	return u.Id
}

// GetUserID user
func (u *UserClaims) GetUserID() string {
	return u.Subject
}

// GetUserName name
func (u *UserClaims) GetUserName() string {
	return u.Name
}

// GetUserRole role
func (u *UserClaims) GetUserRole() string {
	// if u.Role == "" {
	// 	if u.Subject == "1" {
	// 		return "admin" // 作为默认系统用户
	// 	}
	// 	return "default" // 当用户没有权限的时候,启动默认权限
	// }
	return u.Role
}

// SetUserRole role
func (u *UserClaims) SetUserRole(nrole string) string {
	orole := u.Role
	u.Role = nrole
	return orole
}

// GetXidxID user index
func (u *UserClaims) GetXidxID() string {
	return u.XidxID
}

// GetAccountID token
func (u *UserClaims) GetAccountID() string {
	return u.AccountID
}

// GetT3rdID 3rd id
func (u *UserClaims) GetT3rdID() string {
	return u.T3rdID
}

// GetClientID client id
func (u *UserClaims) GetClientID() string {
	return u.ClientID
}

// GetDomain domain
func (u *UserClaims) GetDomain() string {
	return u.Domain
}

// GetIssuer issuer
func (u *UserClaims) GetIssuer() string {
	return u.Issuer
}

// GetAudience audience
func (u *UserClaims) GetAudience() string {
	return u.Audience
}

// GetOrgCode org code
func (u *UserClaims) GetOrgCode() string {
	return u.OrgCode
}

// GetOrgRole org role
func (u *UserClaims) GetOrgRole() string {
	return u.OrgRole
}

// GetOrgDomain org domain
func (u *UserClaims) GetOrgDomain() string {
	return u.OrgDomain
}

// GetOrgAdmin org admin
func (u *UserClaims) GetOrgAdmin() string {
	return u.OrgAdmin
}
