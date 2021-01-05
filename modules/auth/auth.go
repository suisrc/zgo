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
	//helper.UserInfo

	// GetUserName 用户名
	GetUserName() string

	// GetUserID 用户ID
	GetUserID() string
	// GetRoleID 角色ID
	GetRoleID() string
	// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
	GetTokenID() string
	// GetAccountID 登陆ID, 本身不具备任何意义,只是标记登陆方式
	GetAccountID() string

	// 赋予用户临时角色,用户替换,返回之前的角色
	SetRoleID(string) string
	// GetProps() 获取私有属性,该内容会被加密, 注意:内容敏感,不要存储太多的内容
	GetProps() (interface{}, bool)

	// 令牌签发者
	GetIssuer() string
	// 令牌接收者
	GetAudience() string

	// XID
	GetXID() string
	// TID
	GetTID() string
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
