package oauth2

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/store"
)

// Selector 选择器
type Selector map[string]Handler

// NewSelector 全局缓存
func NewSelector(GPA gpa.GPA, Storer store.Storer) (Selector, error) {
	selector := make(map[string]Handler)

	selector["WX"] = &WeixinQm{
		GPA:    GPA,
		Storer: Storer,
	}
	selector["WXQ"] = &WeixinQy{
		GPA:    GPA,
		Storer: Storer,
	}

	return selector, nil
}

//===================================================================================================
//===================================================================================================
//===================================================================================================

// Handler 认证接口
type Handler interface {
	// 获取用户信息
	GetUser(c *gin.Context, rp RequestPlatform, rt RequestToken, relation, userid string) (*UserInfo, error)

	// Handle 处理OAuth2认证, 被动函数回调得到用户的openid
	// func (member bool, openid string, name string) (accountid int, err error)
	Handle(*gin.Context, RequestParams, RequestPlatform, RequestToken, RequestOAuth2) error
}

// UserInfo 目前对接微信， 展示保留这些内容
type UserInfo struct {
	Relation    string // 成员关系
	Name        string // 成员名称
	Mobile      string // 手机号码
	Gender      string // 性别: 0表示未定义, 1表示男性, 2表示女性
	Email       string // 邮箱
	Avatar      string // 头像url
	ThumbAvatar string // 头像缩略图url
	Extattr     map[string]interface{}
}

// RequestParams ...
type RequestParams interface {
	GetCode() string
	GetState() string
	GetScope() string
}

// RequestPlatform ...
type RequestPlatform interface {
	GetID() int64
	GetAppID() string
	GetAppSecret() string
	GetAgentID() string
	GetAgentSecret() string
}

// RequestToken ...
type RequestToken interface {
	// 获取应用访问令牌（从数据库获取）
	FindAccessToken(sqlx *sqlx.DB, token *AccessToken, platform int64) error
	// 存储应用访问令牌（向数据库存储）
	SaveAccessToken(sqlx *sqlx.DB, token *AccessToken) error
	// 异步锁定, 防止多次更新
	LockAsync(sqlx *sqlx.DB, platform int64) error
}

// RequestOAuth2 ...
type RequestOAuth2 interface {
	// 获取重定向Host
	GetRedirectHost() string
	// relation: none (无), member(成员), external(外部联系人)
	FindAccount(relation, openid, userid, deviceid string) (int64, error)
}

//===================================================================================================
//===================================================================================================
//===================================================================================================

// GetRedirectURIByOAuth2Platform Redirect URI by OAuth2Platform
func GetRedirectURIByOAuth2Platform(c *gin.Context, uri string, platform RequestPlatform, oauth2 RequestOAuth2) string {
	if uri == "" {
		uri := oauth2.GetRedirectHost()
		if uri == "" {
			uri = "https://" + c.Request.Host
		}
		uri += c.Request.RequestURI
		uri = url.QueryEscape(uri) // 进行URL编码
	} else if uri[:4] != "http" {
		uri = url.PathEscape("https://"+c.Request.Host) + uri
	}
	return uri
}

// AccessToken xxx
type AccessToken struct {
	Account      int64
	Platform     int64
	TokenID      sql.NullString
	AccessToken  sql.NullString
	ExpiresIn    sql.NullInt64
	ExpiresAt    sql.NullTime
	RefreshToken sql.NullString
	RefreshExpAt sql.NullTime
	Scope        sql.NullString
	AsyncLock    sql.NullTime
}

// ExecWithTokenRetry  exec
// 系统只会重新获取一次令牌
// fn: 需要执行的业务代码, bool -> true: 令牌过期, 获取令牌重试
// ac: 获取token, bool -> true: 必须重新获取令牌
func ExecWithTokenRetry(retry int, fn func(token string) (bool, interface{}, error), ac func(must bool) (string, error)) (interface{}, error) {
	// try := 1
	nac := false
mark:
	token, err := ac(nac)
	if err != nil {
		return nil, err
	} else if token == "" {
		return nil, errors.New("token empty")
	}
	var res interface{}
	if nac, res, err = fn(token); err != nil {
		return res, err
	} else if nac && retry > 0 {
		retry-- //  执行重试
		goto mark
	}
	return res, nil
}

// ExecWithAccessToken access token
func ExecWithAccessToken(c context.Context, fn func(token string) (bool, interface{}, error),
	FindToken func(context.Context) (string, error), NewToken func(context.Context) (string, error),
) (interface{}, error) {
	return ExecWithTokenRetry(1, fn, func(must bool) (string, error) {
		if !must {
			if t, err := FindToken(c); err != nil {
				return "", err
			} else if t != "" {
				return t, nil
			}
		}
		return NewToken(c)
	})
}

// ExecWithAccessTokenX access token
func ExecWithAccessTokenX(c context.Context, fn func(token string) (bool, interface{}, error), h TokenHandler) (interface{}, error) {
	return ExecWithAccessToken(c, fn, h.FindToken, h.NewToken)
}

// TokenHandler t h
type TokenHandler interface {
	FindToken(context.Context) (string, error) // 查询令牌缓存
	NewToken(context.Context) (string, error)  // 新获取令牌
}

//===================================================================================================
//===================================================================================================
//===================================================================================================

var _ TokenHandler = (*TokenManager)(nil)

// TokenManager manager
type TokenManager struct {
	gpa.GPA                                   // 数据库操作
	TokenKey     string                       // 令牌Key, 必须全局唯一
	Platform     int64                        // 平台Kid
	Storer       store.Storer                 // 缓存控制器
	OAuth2Handle RequestToken                 // 访问控制器
	NewTokenFunc func() (*AccessToken, error) // 获取新令牌
	MaxCacheIdle int                          // 使用缓存的临界值, 达到临界值会被动更新令牌, 如果Token是2个小时,推荐使用300秒
	MinCacheIdle int                          // 使用缓存的TTL值, 不要高于MaxCacheIdle,推荐使用60秒
	NewCacheIdle int                          // 有效期低于阈值时候， 异步更新, 推荐使用1800秒
}

// FindToken find
func (a *TokenManager) FindToken(c context.Context) (string, error) {
	value, _, _ := a.Storer.Get(c, a.TokenKey)
	if value != "" {
		// 缓存中存在
		return value, nil
	}
	if a.TokenKey != "" {
		toa2 := &AccessToken{}
		if err := a.OAuth2Handle.FindAccessToken(a.Sqlx, toa2, a.Platform); err != nil {
			if !sqlxc.IsNotFound(err) {
				return "", err
			}
			// 数据不存在
		} else if toa2.Platform == 0 {
			// 数据不存在
		} else if idle := toa2.ExpiresAt.Time.Sub(time.Now()); idle > 0 {
			// 数据存在, 令牌有效
			if a.MaxCacheIdle > 0 && a.MinCacheIdle > 0 {
				// 异步更新
				if idle > time.Duration(a.MaxCacheIdle)*time.Second {
					// 有效期还有300s, 令牌完全可用，将令牌缓存到缓存池中缓存
					a.Storer.Set(c, a.TokenKey, toa2.AccessToken.String, time.Duration(a.MinCacheIdle)*time.Second)
				}
				if idle < time.Duration(a.NewCacheIdle) && toa2.AsyncLock.Valid && toa2.AsyncLock.Time.Add(5*time.Second).Before(time.Now()) {
					// 有效期小于阈值, 异步更新， 锁定的是未来5s中的数据
					if err := a.OAuth2Handle.LockAsync(a.Sqlx, a.Platform); err == nil {
						go a.NewToken(c)
					}
				}
			}
			return toa2.AccessToken.String, nil
		}
	}
	return "", nil // 没有找到,直接返回空字符串就好, error返回非空,会导致程序直接结束
}

// NewToken new
func (a *TokenManager) NewToken(c context.Context) (string, error) {
	token, err := a.NewTokenFunc()
	if err != nil {
		return "", err
	} else if token.AccessToken.String == "" {
		return "", errors.New("access token empty")
	}
	token.AsyncLock = sql.NullTime{Valid: true, Time: time.Now()} // 具有解除锁定的功能
	if err := a.OAuth2Handle.SaveAccessToken(a.Sqlx, token); err != nil {
		return "", err
	}
	a.Storer.Set(c, a.TokenKey, token.AccessToken.String, time.Duration(a.MinCacheIdle)*time.Second)
	return token.AccessToken.String, nil
}

//===================================================================================================

// RequestOAuth2X ...
type RequestOAuth2X struct {
	FindHost func() string
	FindUser func(relation, openid, userid, deviceid string) (int64, error)
}

// GetRedirectHost ...
func (a *RequestOAuth2X) GetRedirectHost() string {
	if a.FindHost == nil {
		return ""
	}
	return a.FindHost()
}

// FindAccount ...
func (a *RequestOAuth2X) FindAccount(relation, openid, userid, deviceid string) (int64, error) {
	return a.FindUser(relation, openid, userid, deviceid)
}

// RequestToken1X ...
type RequestToken1X struct {
	FindToken func(sqlx *sqlx.DB, token *AccessToken, platform int64) error
	SaveToken func(sqlx *sqlx.DB, token *AccessToken) error
	LockToken func(sqlx *sqlx.DB, platform int64) error
}

// FindAccessToken ...
func (a *RequestToken1X) FindAccessToken(sqlx *sqlx.DB, token *AccessToken, platform int64) error {
	if a.FindToken == nil {
		return nil
	}
	return a.FindToken(sqlx, token, platform)
}

// SaveAccessToken ...
func (a *RequestToken1X) SaveAccessToken(sqlx *sqlx.DB, token *AccessToken) error {
	if a.SaveToken == nil {
		return nil
	}
	return a.SaveToken(sqlx, token)
}

// LockAsync ...
func (a *RequestToken1X) LockAsync(sqlx *sqlx.DB, platform int64) error {
	if a.LockToken == nil {
		return nil
	}
	return a.LockToken(sqlx, platform)
}
