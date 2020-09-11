package oauth2

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
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

// Handler 认证接口
type Handler interface {
	// Handle 处理OAuth2认证
	// account 用户账户
	// domain 请求域, 如果不存在,直接指定"", 其作用是在多应用授权时候,准确定位子应用
	// client 请求端, 定位子应用
	Handle(*gin.Context, *schema.SigninOfOAuth2, *schema.SigninGpaOAuth2Platfrm, *schema.SigninGpaAccount) error
}

//===================================================================================================AccessToken-START

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

var _ TokenHandler = (*TokenManager)(nil)

// TokenManager manager
type TokenManager struct {
	gpa.GPA                                                              // 数据库操作
	Key          string                                                  // 令牌Key, 注意不能未空,且必须全局唯一
	PlatformID   int                                                     // 平台ID
	Storer       store.Storer                                            // 缓存
	GetNewToken  func(context.Context, int) (*schema.TokenOAuth2, error) // 获取新令牌
	MaxCacheIdle int                                                     // 使用缓存的临界值, 达到临界值会被动更新令牌, 如果Token是2个小时,推荐使用300秒
	MinCacheIdle int                                                     // 使用缓存的TTL值, 不要高于MaxCacheIdle,推荐使用60秒
}

// FindToken find
func (a *TokenManager) FindToken(c context.Context) (string, error) {
	value, _, _ := a.Storer.Get(c, "access_token:"+a.Key)
	if value != "" {
		// 缓存中存在
		return value, nil
	}
	if a.PlatformID > 0 {
		toa2 := &schema.TokenOAuth2{}
		if err := toa2.QueryByPlatformMust(a.Sqlx, a.PlatformID); err != nil {
			if !sqlxc.IsNotFound(err) {
				return "", err
			}
		}
		if toa2.ID > 0 {
			// 数据库中存在
			if a.MaxCacheIdle > 0 && a.MinCacheIdle > 0 {
				if toa2.ExpiresTime.Time.Sub(time.Now()) > time.Duration(a.MaxCacheIdle)*time.Second {
					// 将令牌缓存到缓存池中缓存
					a.Storer.Set(c, "access_token:"+a.Key, toa2.AccessToken.String, time.Duration(a.MinCacheIdle)*time.Second)
				} else if !toa2.SyncLock.Valid || toa2.SyncLock.Time.Before(time.Now()) {
					// 异步更新令牌, 需要判定锁, 锁定延迟5s, 不主动释放锁定内容
					if err := toa2.LockSync(a.Sqlx); err == nil {
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
	token, err := a.GetNewToken(c, a.PlatformID)
	if err != nil {
		return "", err
	} else if token.AccessToken.String == "" {
		return "", errors.New("access token empty")
	}
	token.SyncLock = sql.NullTime{Valid: true, Time: time.Now()} // 具有解除锁定的功能
	if err := token.UpdateTokenOAuth2(a.Sqlx); err != nil {
		return "", err
	}
	a.Storer.Set(c, "access_token:"+a.Key, token.AccessToken.String, time.Duration(a.MinCacheIdle)*time.Second)
	return token.AccessToken.String, nil
}

//===================================================================================================AccessToken-END
