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

// QueryByID sql select
func (a *SigninGpaUser) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name, status from {{TP}}user where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// SigninGpaRole role
type SigninGpaRole struct {
	ID   int    `db:"id" json:"-"`
	KID  string `db:"kid" json:"id"`
	Name string `db:"name" json:"name"`
}

// QueryByID sql select
func (a *SigninGpaRole) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name from {{TP}}role where id=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByKID sql select
func (a *SigninGpaRole) QueryByKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select id, kid, name from {{TP}}role where kid=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByUserID sql select
func (a *SigninGpaRole) QueryByUserID(sqlx *sqlx.DB, dest *[]SigninGpaRole, userid int) error {
	SQL := "select r.id, r.kid, r.name from {{TP}}user_role ur inner join {{TP}}role r on r.id=ur.role_id where ur.user_id=? and r.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL, userid)
}

// SigninGpaAccount account
type SigninGpaAccount struct {
	ID           int            `db:"id"`
	PID          sql.NullInt64  `db:"pid"`
	Account      string         `db:"account"`
	AccountType  int            `db:"account_typ"`
	AccountKind  sql.NullInt64  `db:"account_kid"`
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
	sqr := strings.Builder{}
	sqr.WriteString("select id")
	sqr.WriteString(", pid")
	sqr.WriteString(", account")
	sqr.WriteString(", account_typ")
	sqr.WriteString(", account_kid")
	sqr.WriteString(", password")
	sqr.WriteString(", password_salt")
	sqr.WriteString(", password_type")
	sqr.WriteString(", verify_type")
	sqr.WriteString(", verify_secret")
	sqr.WriteString(", user_id")
	sqr.WriteString(", role_id")
	sqr.WriteString(" from {{TP}}account")
	sqr.WriteString(" where account=? and account_typ=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		sqr.WriteString(" and account_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and account_kid is null")
	}
	sqr.WriteString(" and status=1")
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// SigninGpaOAuth2Account account
type SigninGpaOAuth2Account struct {
	KID string `db:"kid"`
}

// QueryByUdx kid
func (a *SigninGpaOAuth2Account) QueryByUdx(sqlx *sqlx.DB, accountID, clientID int, userKID, roleKID string) error {
	SQL := "select kid from {{TP}}oauth2_account where account_id=? and client_id=? and user_kid=? and role_kid=?"
	return sqlx.Get(a, SQL, accountID, clientID, userKID, roleKID)
}
