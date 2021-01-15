package jwt

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/suisrc/zgo/modules/auth"
)

var _ auth.TokenInfo = &TokenInfo{}

// TokenInfo 令牌信息
type TokenInfo struct {
	// TokenType    string `json:"token_type,omitempty"`    // 令牌类型
	TokenID      string `json:"token_id,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`    // 访问令牌
	ExpiresAt    int64  `json:"expires_at,omitempty"`      // 访问令牌过期时间
	RefreshToken string `json:"refresh_token,omitempty"`   // 刷新令牌
	RefreshExpAt   int64  `json:"refresh_expires,omitempty"` // 刷新令牌过期时间
}

// GetTokenID token id
func (t *TokenInfo) GetTokenID() string {
	return t.TokenID
}

// GetAccessToken access token
func (t *TokenInfo) GetAccessToken() string {
	return t.AccessToken
}

// GetExpiresAt expires at
func (t *TokenInfo) GetExpiresAt() int64 {
	return t.ExpiresAt
}

// GetRefreshToken refresh token
func (t *TokenInfo) GetRefreshToken() string {
	return t.RefreshToken
}

// GetRefreshExpAt refresh expires
func (t *TokenInfo) GetRefreshExpAt() int64 {
	return t.RefreshExpAt
}

// EncodeToJSON to json
func (t *TokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}

//=================================================
// 分割线
//=================================================

// GetBearerToken 获取用户令牌
func GetBearerToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		prefix := "Bearer "
		if auth := c.GetHeader("Authorization"); auth != "" && strings.HasPrefix(auth, prefix) {
			return auth[len(prefix):], nil
		}
	}

	return "", auth.ErrNoneToken
}

// GetQueryToken 获取用户令牌
func GetQueryToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if auth, ok := c.GetQuery("token"); ok && auth != "" {
			return auth, nil
		}
	}

	return "", auth.ErrNoneToken
}

// GetCookieToken 获取用户令牌
func GetCookieToken(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if auth, err := c.Cookie("authorization"); err == nil && auth != "" {
			return auth, nil
		}
	}

	return "", auth.ErrNoneToken
}
