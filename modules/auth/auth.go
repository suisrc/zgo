package auth

import (
	"context"
	"errors"
)

// 定义错误
var (
	// ErrInvalidToken 无效令牌
	ErrInvalidToken = errors.New("invalid token")
	// ErrNoneToken 没有令牌
	ErrNoneToken = errors.New("none token")
)

// TokenInfo 令牌信息
type TokenInfo interface {

	// 获取访问令牌
	GetAccessToken() string

	// 获取令牌到期时间戳
	GetExpiresAt() int64

	// 获取刷新令牌
	GetRefreshToken() string

	// 获取刷新令牌过期时间戳
	GetRefreshExp() int64

	// JSON
	EncodeToJSON() ([]byte, error)
}

// UserInfo user
type UserInfo interface {
	// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
	GetTokenID() string

	// GetUserID 用户ID
	GetUserID() string
	// GetUserName 用户名
	GetUserName() string
	// GetUserRole 角色ID
	GetUserRole() string
	// GetXidxID 直接获取用户索引ID
	GetXidxID() string
	// GetUserAccount 登陆ID, 本身不具备任何意义,只是标记登陆方式
	GetAccountID() string
	// GetT3rdID 获取用户第三方索引
	GetT3rdID() string
	// GetT3rdID 获取用户第三方索引
	GetClientID() string

	// GetDomain
	GetDomain() string
	// GetIssuer 令牌签发者
	GetIssuer() string
	// GetAudience 令牌接收者
	GetAudience() string

	// GetOrgCode
	GetOrgCode() string
	// GetOrgRole
	GetOrgRole() string
	// GetOrgDomain
	GetOrgDomain() string
	// IsOrgAdmin
	GetOrgAdmin() string

	// 赋予用户临时角色,用户替换,返回之前的角色
	SetUserRole(string) string
}

// Auther 认证接口
type Auther interface {
	// GetUserInfo 获取用户
	GetUserInfo(c context.Context) (UserInfo, error)

	// GenerateToken 生成令牌
	GenerateToken(c context.Context, u UserInfo) (TokenInfo, UserInfo, error)

	// RefreshToken 刷新令牌
	RefreshToken(c context.Context, t string, f func(UserInfo, int) error) (TokenInfo, UserInfo, error)

	// DestroyToken 销毁令牌
	DestroyToken(c context.Context, u UserInfo) error

	// UpdateAuther 更新
	UpdateAuther(c context.Context) error
}
