package service

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
	"github.com/suisrc/zgo/modules/store/buntdb"
)

// AuthOpts 认证配置信息
// 认证需要频繁操作,所以这里需要使用内存缓存
type AuthOpts struct {
	gpa.GPA
	Store store.Storer
	Jwts  map[interface{}]*schema.JwtGpaOpts
}

// NewStorer 全局缓存
func NewStorer() (store.Storer, func(), error) {
	store, err := buntdb.NewStore(":memory:") // 使用内存缓存
	if err != nil {
		return nil, nil, err
	}
	return store, func() { store.Close() }, nil
}

// NewAuther of auth.Auther
// 授权认证使用的auther内容
func NewAuther(opts *AuthOpts) auth.Auther {
	secret := config.C.JWTAuth.SigningSecret
	if secret == "" {
		secret = crypto.UUID(128)
		logger.Infof(nil, "jwt secret: %s", secret)
	}

	opts.Jwts = map[interface{}]*schema.JwtGpaOpts{}
	auther := jwt.New(opts.Store,
		jwt.SetKeyFunc(opts.keyFunc),
		jwt.SetNewClaims(opts.signingFunc),
		jwt.SetFixClaimsFunc(opts.claimsFunc),
		jwt.SetUpdateFunc(opts.updateFunc),
		jwt.SetTokenFunc(opts.tokenFunc),
		jwt.SetSigningSecret(secret),                  // 注册令牌签名密钥
		jwt.SetExpired(config.C.JWTAuth.LimitExpired), // 访问令牌生命周期
		jwt.SetRefresh(config.C.JWTAuth.LimitRefresh), // 刷新令牌声明周期
	)
	// 触发updateFunc方法
	go auther.UpdateAuther(nil)
	return auther
}

// 获取令牌的方式
func (a *AuthOpts) tokenFunc(ctx context.Context) (string, error) {
	if c, ok := ctx.(*gin.Context); ok {
		if state := c.Query("state"); state != "" {
			// 使用state隐藏令牌,一般用于重定向回调上
			if token, _, _ := a.Store.Get(ctx, state); token != "" {
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

// 更新认证
func (a *AuthOpts) updateFunc(c context.Context) error {
	opts := map[interface{}]*schema.JwtGpaOpts{}
	jwtOpt := new(schema.JwtGpaOpts)
	jwtOpts := []schema.JwtGpaOpts{}
	if err := jwtOpt.QueryAll(a.Sqlx, &jwtOpts); err != nil {
		logger.Errorf(c, logger.ErrorWW(err)) // 更新发生异常
	} else {
		for _, v := range jwtOpts {
			v.SecretByte = []byte(v.Secret)
			opts[v.KID] = &v
		}
	}
	a.Jwts = opts

	return nil
}

// 修正令牌
func (a *AuthOpts) claimsFunc(c context.Context, claims *jwt.UserClaims) error {
	if kid, ok := helper.GetCtxValueToString(c, helper.ResJwtKey); ok {
		opt, ok := a.Jwts[kid]
		if !ok {
			return errors.New("signing jwt, kid error")
		}
		if opt.Expired > 0 {
			now := time.Unix(claims.IssuedAt, 0)
			claims.ExpiresAt = now.Add(time.Duration(opt.Expired) * time.Second).Unix() // 修改时间
		}
		if opt.Audience.Valid {
			claims.Audience = opt.Audience.String
		}
		if opt.Issuer.Valid {
			claims.Issuer = opt.Issuer.String
		}
	}
	return nil
}

// 获取jwt密钥
func (a *AuthOpts) keyFunc(c context.Context, token *jwtgo.Token, method jwtgo.SigningMethod, secret interface{}) (interface{}, error) {
	token.Method = method // 强制使用配置, 防止alg使用none而跳过验证

	// 获取处理的密钥
	if kid, ok := token.Header["kid"]; ok {
		helper.SetCtxValue(c, helper.ResJwtKey, kid)
		if opt, ok := a.Jwts[kid]; ok {
			return opt.SecretByte, nil
		}
		return nil, errors.New("parse jwt, kid error")
	}
	return secret, nil // 使用默认令牌，默认方法进行解密
	// return nil, auth.ErrInvalidToken // 无法验证
}

// 签名jwt令牌
func (a *AuthOpts) signingFunc(c context.Context, claims *jwt.UserClaims, method jwtgo.SigningMethod, secret interface{}) (string, error) {
	if kid, ok := helper.GetCtxValueToString(c, helper.ResJwtKey); ok {
		// 使用jwt私有密钥
		if opt, ok := a.Jwts[kid]; ok {
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
