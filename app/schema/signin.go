package schema

import (
	"database/sql"

	"github.com/suisrc/zgo/modules/auth"
)

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password" binding:"required"` // 密码
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	Role     string `json:"role"`                        // 角色
	Attach   string `json:"attach"`                      // reset:重置登入 refresh:刷新令牌
	Client   string `json:"client"`                      // 应用, 默认为主应用, 为空
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
	UserName string
	UserID   string
	RoleID   string
	TokenID  string
	Issuer   string
	Audience string
	SIID     string
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
	return s.SIID
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
	UID    string `db:"uid" json:"id"`
	Name   string `db:"name" json:"name"`
	Status bool   `db:"status" json:"-"`
}

// SQLByID sql select
func (*SigninGpaUser) SQLByID() string {
	return "select id, uid, name, status from user where id=?"
}

// SigninGpaRole role
type SigninGpaRole struct {
	ID   int    `db:"id" json:"-"`
	UID  string `db:"uid" json:"id"`
	Name string `db:"name" json:"name"`
}

// SQLByID sql select
func (*SigninGpaRole) SQLByID() string {
	return "select id, uid, name from  role where id=? and status=1"
}

// SQLByUID sql select
func (*SigninGpaRole) SQLByUID() string {
	return "select id, uid, name from role where uid=? and status=1"
}

// SQLByName sql select
func (*SigninGpaRole) SQLByName() string {
	return "select id, uid, name from role where name=? and status=1"
}

// SQLByUserID sql select
func (*SigninGpaRole) SQLByUserID() string {
	return "select r.id, r.uid, r.name from user_role ur inner join role r on r.id=ur.role_id where ur.user_id=? and r.status=1"
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
	VerifyType   sql.NullString `db:"verify_type"`
	Password     sql.NullString `db:"password"`
	PasswordSalt sql.NullString `db:"password_salt"`
	PasswordType sql.NullString `db:"password_type"`
	UserID       int            `db:"user_id"`
	RoleID       sql.NullInt64  `db:"role_id"`
	OAuth2ID     sql.NullInt64  `db:"oauth2_id"`

	// SQLX1 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
	// SQLX2 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
}

// SQLByAccount sql select
func (*SigninGpaAccount) SQLByAccount() string {
	return "select id, verify_type, password, password_salt, password_type, user_id, role_id, oauth2_id from account where account=? and account_type='user' and platform='ZGO' and status=1"
}
