package injector

import (
	"context"
	"errors"
	"time"

	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/auth/jwt/store/buntdb"
)

// AuthOpts 认证配置信息
// 认证需要频繁操作,所以这里需要使用内存缓存
type AuthOpts struct {
	service.GPA
	jwts map[interface{}]*schema.JwtGpaOpts
}

// NewAuther of auth.Auther
// 授权认证使用的auther内容
func NewAuther(opts *AuthOpts) auth.Auther {
	store, err := buntdb.NewStore(":memory:") // 使用内存缓存
	if err != nil {
		panic(err)
	}
	secret := config.C.JWTAuth.SigningSecret
	if secret == "" {
		secret = crypto.UUID(128)
		logger.Infof(nil, "jwt secret: %s", secret)
	}

	opts.jwts = map[interface{}]*schema.JwtGpaOpts{}
	auther := jwt.New(store,
		jwt.SetSigningSecret(secret), // 注册令牌签名密钥
		jwt.SetKeyFunc(opts.keyFunc),
		jwt.SetNewClaims(opts.signingFunc),
		jwt.SetFixClaimsFunc(opts.claimsFunc),
		jwt.SetUpdateFunc(opts.updateFunc),
	)
	// 触发updateFunc方法
	go auther.UpdateAuther(nil)
	return auther
}

// 更新认证
func (a *AuthOpts) updateFunc(c context.Context) error {
	opts := map[interface{}]*schema.JwtGpaOpts{}
	jwtOpt := new(schema.JwtGpaOpts)
	jwtOpts := []schema.JwtGpaOpts{}
	if err := jwtOpt.QueryAll(a.Sqlx, jwtOpts); err != nil {
		logger.Errorf(c, err.Error()) // 更新发生异常
	} else {
		for _, v := range jwtOpts {
			opts[v.KID] = &v
		}
	}
	a.jwts = opts

	return nil
}

// 修正令牌
func (a *AuthOpts) claimsFunc(c context.Context, claims *jwt.UserClaims) error {
	if kid, ok := helper.GetJwtKid(c); ok {
		opt, ok := a.jwts[kid]
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
		helper.SetJwtKid(c, kid)
		if opt, ok := a.jwts[kid]; ok {
			return opt.Secret, nil
		}
		return nil, errors.New("parse jwt, kid error")
	}
	return secret, nil // 使用默认令牌，默认方法进行解密
	// return nil, auth.ErrInvalidToken // 无法验证
}

// 签名jwt令牌
func (a *AuthOpts) signingFunc(c context.Context, claims jwtgo.Claims, method jwtgo.SigningMethod, secret interface{}) (string, error) {
	if kid, ok := helper.GetJwtKid(c); ok {
		// 使用jwt私有密钥
		if opt, ok := a.jwts[kid]; ok {
			token := &jwtgo.Token{
				Header: map[string]interface{}{
					"typ": "JWT",
					"alg": method.Alg(),
					"kid": kid,
				},
				Claims: claims,
				Method: method,
			}
			return token.SignedString(opt.Secret)
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
