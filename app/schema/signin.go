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

// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
func (s *SigninUser) GetTokenID() string {
	return s.TokenID
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

// UserSignin user
type UserSignin struct {
	ID   int    `db:"id" json:"-"`
	UID  string `db:"uid" json:"id"`
	Name string `db:"name" json:"name"`
}

// RoleSignin role
type RoleSignin struct {
	ID   int    `db:"id" json:"-"`
	UID  string `db:"uid" json:"id"`
	Name string `db:"name" json:"name"`
}

// ClientSignin client
type ClientSignin struct {
	ID       int            `db:"id"`
	Issuer   sql.NullString `db:"issuer"`
	Audience sql.NullString `db:"audience"`
}

// AccountSignin account
type AccountSignin struct {
	ID           int            `db:"id"`
	Account      string         `db:"account"`
	AccountType  string         `db:"account_type"`
	Platform     string         `db:"platform"`
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
