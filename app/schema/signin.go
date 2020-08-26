package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/modules/auth"
)

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password" binding:"required"` // 密码
	KID      string `json:"kid"`                         // 授权平台
	Client   string `json:"client"`                      // 子应用ID
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	Role     string `json:"role"`                        // 角色
}

// SigninResult 登陆返回值
type SigninResult struct {
	Status       string        `json:"status" default:"ok"`    // 'ok' | 'error' 不适用boolean类型是为了以后可以增加扩展
	Token        string        `json:"token,omitempty"`        // 令牌
	Expired      int64         `json:"expired,omitempty"`      // 过期时间
	RefreshToken string        `json:"refreshToken,omitempty"` // 刷新令牌
	Message      string        `json:"message,omitempty"`      // 消息,有限显示
	Roles        []interface{} `json:"roles,omitempty"`        // 多角色的时候，返回角色，重新确认登录
}

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	UserName  string
	UserID    string
	RoleID    string
	TokenID   string
	Issuer    string
	Audience  string
	AccountID string
}

// GetUserName 用户名
func (s *SigninUser) GetUserName() string {
	return s.UserName
}

// GetUserID 用户ID
func (s *SigninUser) GetUserID() string {
	return s.UserID
}

// GetRoleID 角色ID
func (s *SigninUser) GetRoleID() string {
	return s.RoleID
}

// SetRoleID 角色ID
func (s *SigninUser) SetRoleID(nrole string) string {
	orole := s.RoleID
	s.RoleID = nrole
	return orole
}

// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
func (s *SigninUser) GetTokenID() string {
	return s.TokenID
}

// GetAccountID token
func (s *SigninUser) GetAccountID() string {
	return s.AccountID
}

// GetProps 获取私有属性,该内容会被加密, 注意:内容敏感,不要存储太多的内容
func (s *SigninUser) GetProps() (interface{}, bool) {
	return nil, false
}

// GetIssuer 令牌签发者
func (s *SigninUser) GetIssuer() string {
	return s.Issuer
}

// GetAudience 令牌接收者
func (s *SigninUser) GetAudience() string {
	return s.Audience
}

//==============================================================================

// SigninGpaUser user
type SigninGpaUser struct {
	ID     int    `db:"id" json:"-"`
	KID    string `db:"kid" json:"id"`
	Name   string `db:"name" json:"name"`
	Status bool   `db:"status" json:"-"`
}

// SQLByID sql select
func (*SigninGpaUser) SQLByID() string {
	return "select id, kid, name, status from user where id=?"
}

// SigninGpaRole role
type SigninGpaRole struct {
	ID   int    `db:"id" json:"-"`
	KID  string `db:"kid" json:"id"`
	Name string `db:"name" json:"name"`
}

// SQLByID sql select
func (*SigninGpaRole) SQLByID() string {
	return "select id, kid, name from  role where id=? and status=1"
}

// SQLByKID sql select
func (*SigninGpaRole) SQLByKID() string {
	return "select id, kid, name from role where kid=? and status=1"
}

// SQLByName sql select
func (*SigninGpaRole) SQLByName() string {
	return "select id, kid, name from role where name=? and status=1"
}

// SQLByUserID sql select
func (*SigninGpaRole) SQLByUserID() string {
	return "select r.id, r.kid, r.name from user_role ur inner join role r on r.id=ur.role_id where ur.user_id=? and r.status=1"
}

// SigninGpaClient client
type SigninGpaClient struct {
	ID       int            `db:"id"`
	Issuer   sql.NullString `db:"issuer"`
	Audience sql.NullString `db:"audience"`
}

// SQLByClientKey sql select
func (*SigninGpaClient) SQLByClientKey() string {
	return "select id, issuer, audience from user where client_key=?"
}

// SigninGpaAccount account
type SigninGpaAccount struct {
	ID           int            `db:"id"`
	PID          sql.NullInt32  `db:"pid"`
	Account      string         `db:"account"`
	AccountType  int            `db:"account_typ"`
	AccountKind  sql.NullInt32  `db:"account_kid"`
	Password     sql.NullString `db:"password"`
	PasswordSalt sql.NullString `db:"password_salt"`
	PasswordType sql.NullString `db:"password_type"`
	VerifyType   sql.NullString `db:"verify_type"`
	VerifySecret sql.NullString `db:"verify_secret"`
	UserID       int            `db:"user_id"`
	RoleID       sql.NullInt64  `db:"role_id"`

	// SQLX1 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
	// SQLX2 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
}

// QueryByAccount sql select
func (a *SigninGpaAccount) QueryByAccount(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	SQL := strings.Builder{}
	SQL.WriteString("select id")
	SQL.WriteString(", pid")
	SQL.WriteString(", account")
	SQL.WriteString(", account_typ")
	SQL.WriteString(", account_kid")
	SQL.WriteString(", password")
	SQL.WriteString(", password_salt")
	SQL.WriteString(", password_type")
	SQL.WriteString(", verify_type")
	SQL.WriteString(", verify_secret")
	SQL.WriteString(", user_id")
	SQL.WriteString(", role_id")
	SQL.WriteString(" from account")
	SQL.WriteString(" where account=? and account_typ=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		SQL.WriteString(" and account_kid=?")
		params = append(params, kid)
	} else {
		SQL.WriteString(" and account_kid is null")
	}
	SQL.WriteString(" and status=1")
	return sqlx.Get(a, SQL.String(), params...)
}

// SigninGPAOAuth2Account account
type SigninGPAOAuth2Account struct {
	KID string `db:"kid"`
}

// SQLByKID kid
func (*SigninGPAOAuth2Account) SQLByKID() string {
	return "select kid where account_id=? and client_id=? and user_kid=? and role_kid=?"
}
