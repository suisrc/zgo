package helper

import (
	"context"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suisrc/zgo/modules/auth"
)

// 定义上下文中的键
const (
	Prefix      = "zgo"
	UserInfoKey = Prefix + "/user-info"   // user info
	TraceIDKey  = Prefix + "/tract-id"    // trace id
	ReqBodyKey  = Prefix + "/req-body"    // request body
	ResBodyKey  = Prefix + "/res-body"    // response body
	ResJwtKey   = Prefix + "/res-jwt-kid" // jwt kid
	ResTknKey   = Prefix + "/res-tkn-kid" // tkn kid

	XReqOriginHostKey   = "X-Request-Origin-Host"
	XReqOriginPathKey   = "X-Request-Origin-Path"
	XReqOriginMethodKey = "X-Request-Origin-Method"
)

// UserInfoFunc user
type UserInfoFunc interface {
	GetUserInfo() (auth.UserInfo, bool)
	SetUserInfo(auth.UserInfo)
}

// GetUserInfo 用户
func GetUserInfo(c *gin.Context) (auth.UserInfo, bool) {
	if v, ok := c.Get(UserInfoKey); ok {
		if u, b := v.(auth.UserInfo); b {
			return u, true
		}
	}
	return nil, false
}

// SetUserInfo 用户
func SetUserInfo(c *gin.Context, user auth.UserInfo) {
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
	// log.Println(traceID)
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

// GetHostIP 获取主机端IP
func GetHostIP(c *gin.Context) string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String()
				}

			}
		}
	}

	return ""
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

// GetCtxValue 获取令牌加密方式
func GetCtxValue(ctx context.Context, key string) (interface{}, bool) {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Get(key)
	}
	return "", false
}

// GetCtxValueToString 获取令牌加密方式
func GetCtxValueToString(ctx context.Context, key string) (string, bool) {
	if c, ok := ctx.(*gin.Context); ok {
		if v, ok := c.Get(key); ok {
			if s, ok := v.(string); ok {
				return s, true
			}
		}
	}
	return "", false
}

// SetCtxValue 配置令牌加密方式
func SetCtxValue(ctx context.Context, key string, value interface{}) bool {
	if c, ok := ctx.(*gin.Context); ok {
		c.Set(key, value)
		return true
	}
	return false
}

// Now 获取当前时间
// func Now() time.Time {
// 	//return time.Now().In(time.Local)
// 	return time.Now()
// }
