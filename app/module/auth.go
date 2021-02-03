package module

import (
	"context"
	"errors"
	"strings"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"
)

// AuthOpts 认证配置信息
// 认证需要频繁操作,所以这里需要使用内存缓存
type AuthOpts struct {
	gpa.GPA
	Storer         store.Storer
	CachedJwtOtps1 map[interface{}]*AuthJwt // 加密配置
	CachedExpireAt time.Time                // 刷新时间
}

// AuthJwt 验证器
type AuthJwt struct {
	JwtOpt   *schema.ClientGpaWebToken // 加密配置
	ExpireAt time.Time                 // 刷新时间
}

// NewAuther of auth.Auther
// 授权认证使用的auther内容
func NewAuther(opts *AuthOpts) auth.Auther {
	secret := config.C.JWTAuth.SigningSecret
	if secret == "" {
		secret = crypto.UUID(128)
		logger.Infof(nil, "jwt secret: %s", secret)
	}

	opts.CachedJwtOtps1 = map[interface{}]*AuthJwt{}
	opts.CachedExpireAt = time.Now().Add(2 * time.Minute)
	auther := jwt.New(opts.Storer,
		jwt.SetKeyFunc(opts.key),
		jwt.SetNewClaims(opts.signing),
		jwt.SetFixClaimsFunc(opts.claims),
		jwt.SetUpdateFunc(opts.update),
		jwt.SetTokenFunc(opts.token),
		jwt.SetSigningSecret(secret),                  // 注册令牌签名密钥
		jwt.SetExpired(config.C.JWTAuth.LimitExpired), // 访问令牌生命周期
		jwt.SetRefresh(config.C.JWTAuth.LimitRefresh), // 刷新令牌声明周期
	)
	// 触发updateFunc方法
	go auther.UpdateAuther(nil)
	return auther
}

// clear 清理缓存
func (a *AuthOpts) clear(force bool, kid string) {
	if a.CachedJwtOtps1 == nil {
		// do nothing
	} else if force {
		a.CachedJwtOtps1 = nil
		a.CachedExpireAt = time.Now().Add(AuthCachedExpireAt)
	} else if kid != "" {
		delete(a.CachedJwtOtps1, kid) // 清除指定缓存
	} else {
		now := time.Now()
		for k, v := range a.CachedJwtOtps1 {
			if v.ExpireAt.Before(now) {
				delete(a.CachedJwtOtps1, k) // 清除过期缓存
			}
		}
	}
}

// jwt 获取加密认证的jwt配置信息
func (a *AuthOpts) opts(ctx context.Context, kid interface{}) (*schema.ClientGpaWebToken, bool) {
	if a.CachedJwtOtps1 == nil {
		a.CachedJwtOtps1 = map[interface{}]*AuthJwt{}
		a.CachedExpireAt = time.Now().Add(AuthCachedExpireAt)
	} else if a.CachedExpireAt.Before(time.Now()) {
		a.CachedExpireAt = time.Now().Add(AuthCachedExpireAt) // 设定04分钟后刷新
		// a.clear(false, "")
		go a.clear(false, "")
		// defer func() { go a.clear(false, "") }()
	}

	if jwt, ok := a.CachedJwtOtps1[kid]; ok {
		return jwt.JwtOpt, ok // 使用缓存
	}
	if kidstr, ok := kid.(string); ok {
		cgw := schema.ClientGpaWebToken{}
		if err := cgw.QueryByKID(a.Sqlx, kidstr); err == nil && cgw.KID != "" {
			cgw.SecretByte = []byte(cgw.JwtSecret.String)
			a.CachedJwtOtps1[kid] = &AuthJwt{
				JwtOpt:   &cgw,
				ExpireAt: time.Now().Add(2 * AuthCachedExpireAt),
			}
			// log.Println(cgw)
			return &cgw, true // 从数据库获取
		}
	}
	return nil, false
}

// 更新认证
func (a *AuthOpts) update(c context.Context) error {
	// 使用按需加载的方式， 所以这里的更新， 只要清除缓存即可， 我们可以通过jwtopt方法重新加载缓存
	a.clear(true, "")
	return nil
}

// 获取令牌的方式
func (a *AuthOpts) token(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if state := c.Query("state"); state != "" {
			// 使用state隐藏令牌,一般用于重定向回调上
			if token, _, _ := a.Storer.Get(ctx, state); token != "" {
				return token, nil
			}
		}
		prefix := "Bearer "
		if auth := c.GetHeader("Authorization"); auth != "" && strings.HasPrefix(auth, prefix) {
			return auth[len(prefix):], nil
		}
	}
	return "", auth.ErrNoneToken
}

// 修正令牌
func (a *AuthOpts) claims(c context.Context, claims *jwt.UserClaims) (int, error) {
	if kid, ok := helper.GetCtxValueToString(c, helper.ResJwtKey); ok {
		opt, ok := a.opts(c, kid)
		if !ok {
			return -1, errors.New("signing jwt, kid error")
		}
		if opt.JwtExpired.Valid && opt.JwtExpired.Int64 > 0 {
			now := time.Unix(claims.IssuedAt, 0)
			claims.ExpiresAt = now.Add(time.Duration(opt.JwtExpired.Int64) * time.Second).Unix() // 修改时间
		}
		if opt.JwtAudience.Valid {
			claims.Audience = opt.JwtAudience.String
		}
		if opt.JwtIssuer.Valid {
			claims.Issuer = opt.JwtIssuer.String
		}
		if opt.JwtRefresh.Valid && opt.JwtRefresh.Int64 > 0 {
			return int(opt.JwtRefresh.Int64), nil
		}
	}
	return -1, nil
}

// 获取jwt密钥
func (a *AuthOpts) key(c context.Context, token *jwtgo.Token, method jwtgo.SigningMethod, secret interface{}) (interface{}, error) {
	token.Method = method // 强制使用配置, 防止alg使用none而跳过验证

	// 获取处理的密钥
	if kid, ok := token.Header["kid"]; ok {
		helper.SetCtxValue(c, helper.ResJwtKey, kid)
		if opt, ok := a.opts(c, kid); ok {
			return opt.SecretByte, nil
		}
		return nil, errors.New("parse jwt, kid error")
	}
	return secret, nil // 使用默认令牌，默认方法进行解密
	// return nil, auth.ErrInvalidToken // 无法验证
}

// 签名jwt令牌
func (a *AuthOpts) signing(c context.Context, claims *jwt.UserClaims, method jwtgo.SigningMethod, secret interface{}) (string, error) {
	if kid, ok := helper.GetCtxValueToString(c, helper.ResJwtKey); ok {
		// 使用jwt私有密钥
		if opt, ok := a.opts(c, kid); ok {
			token := &jwtgo.Token{
				Header: map[string]interface{}{
					"typ": "JWT",
					"alg": method.Alg(),
					"kid": kid,
				},
				Claims: claims,
				Method: method,
			}
			return token.SignedString(opt.SecretByte)
		}
		return "", errors.New("signing jwt kid error")
	}
	// 使用公共密钥
	token := &jwtgo.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
		},
		Claims: claims,
		Method: method,
	}
	return token.SignedString(secret)
}
