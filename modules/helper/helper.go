package helper

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 定义上下文中的键
const (
	Prefix       = "zgo"
	UserInfoKey  = Prefix + "/user-info"
	TraceIDKey   = Prefix + "/tract-id"
	ReqBodyKey   = Prefix + "/req-body"
	ResBodyKey   = Prefix + "/res-body"
	ResJwtKey    = Prefix + "/res-jwt-kid"
	ResJwtOptKey = Prefix + "/res-jwt-opt"

	XReqUserKey = "X-Request-User-KID"
	XReqRoleKey = "X-Request-Role-KID"
)

// UserInfo 用户信息
type UserInfo interface {
	GetUserID() string
	GetRoleID() string
	GetProps() (interface{}, bool)

	GetUserName() string
	GetTokenID() string
	GetIssuer() string
	GetAudience() string

	GetAccountID() string
	SetRoleID(string) string
}

// UserInfoFunc user
type UserInfoFunc interface {
	GetUserInfo() (UserInfo, bool)
	SetUserInfo(UserInfo)
}

// GetUserInfo 用户
func GetUserInfo(c *gin.Context) (UserInfo, bool) {
	if v, ok := c.Get(UserInfoKey); ok {
		if u, b := v.(UserInfo); b {
			return u, true
		}
	}
	return nil, false
}

// SetUserInfo 用户
func SetUserInfo(c *gin.Context, user UserInfo) {
	c.Set(UserInfoKey, user)
}

// GetTraceID 根据山下问,获取追踪ID
func GetTraceID(c *gin.Context) string {
	if c == nil {
		v, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		return v.String()
	}
	if v, ok := c.Get(TraceIDKey); ok && v != "" {
		return v.(string)
	}

	// 优先从请求头中获取请求ID
	traceID := c.GetHeader("X-Request-Id")
	if traceID == "" {
		// 没有自建
		v, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		traceID = v.String()
	}
	c.Set(TraceIDKey, traceID)
	return traceID
}

// GetClientIP 获取客户端IP
func GetClientIP(c *gin.Context) string {
	if v := c.GetHeader("X-Forwarded-For"); v != "" {
		if len := strings.Index(v, ","); len > 0 {
			return v[:len]
		}
		return v
	}
	return c.ClientIP()
}

// GetAcceptLanguage 获取浏览器语言
func GetAcceptLanguage(c *gin.Context) string {
	if v := c.GetHeader("Accept-Language"); v != "" {
		if len := strings.Index(v, ","); len > 0 {
			v = v[:len]
		}
		if len := strings.Index(v, ";"); len > 0 {
			v = v[:len]
		}
		return v
	}
	return ""
}

// GetJwtKid 获取令牌加密方式
func GetJwtKid(ctx context.Context) (interface{}, bool) {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Get(ResJwtKey)
	}
	return "", false
}

// GetJwtKidStr 获取令牌加密方式
func GetJwtKidStr(ctx context.Context) (string, bool) {
	if c, ok := ctx.(*gin.Context); ok {
		if v, ok := c.Get(ResJwtKey); ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

// SetJwtKid 配置令牌加密方式
func SetJwtKid(ctx context.Context, kid interface{}) bool {
	if c, ok := ctx.(*gin.Context); ok {
		c.Set(ResJwtKey, kid)
		return true
	}
	return false
}

// Now 获取当前时间
// func Now() time.Time {
// 	//return time.Now().In(time.Local)
// 	return time.Now()
// }
