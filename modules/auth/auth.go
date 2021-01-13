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

	// 获取令牌ID
	GetTokenID() string

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
	// GetUserAccount 登陆ID, 本身不具备任何意义,只是标记登陆方式, 使用token反向加密
	GetAccountID() string
	// GetUserIdxID 直接获取用户索引ID, 使用token反向加密
	GetUserIdxID() string

	// GetUserID 用户ID， GetOrgCode不为空(P6M开头的租户除外)，不提供
	GetUserID() string
	// GetUserName 用户名， GetOrgCode不为空(P6M开头的租户除外)，不提供
	GetUserName() string
	// GetUserRoles 角色， GetOrgCode不为空(P6M开头的租户除外)，不提供
	GetUserRoles() []string

	// GetOrgCode
	GetOrgCode() string
	// IsOrgAdmin 'admin'为用户管理员， GetOrgCode为空，提供
	GetOrgAdmin() string

	// GetDomain 领域标识
	GetDomain() string
	// GetIssuer 令牌签发者
	GetIssuer() string
	// GetAudience 令牌接收者
	GetAudience() string

	// 通过服务名称获取当前服务的角色
	GetUserSvcRoles(string) []string
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
