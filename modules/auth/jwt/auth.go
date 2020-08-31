package jwt

/*
 为什么使用反向验证(只记录登出的用户, 因为我们确信点击登出的操作比点击登陆的操作要少的多的多)
*/
import (
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"
)

type options struct {
	tokenType     string                                                                                 // 令牌类型,传递给TokenInfo
	expired       int                                                                                    // 过期间隔
	signingMethod jwt.SigningMethod                                                                      // 签名方法
	signingSecret interface{}                                                                            // 公用签名密钥
	signingFunc   func(context.Context, jwt.Claims, jwt.SigningMethod, interface{}) (string, error)      // JWT构建令牌
	claimsFunc    func(context.Context, *UserClaims) error                                               // 处理令牌
	keyFunc       func(context.Context, *jwt.Token, jwt.SigningMethod, interface{}) (interface{}, error) // JWT中获取密钥, 该内容可以忽略默认的signingMethod和signingSecret
	parseFunc     func(context.Context, string) (*UserClaims, error)                                     // 解析令牌
	tokenFunc     func(context.Context) (string, error)                                                  // 获取令牌
	updateFunc    func(context.Context) error                                                            // 更新Auther
}

// Option 定义参数项
type Option func(*options)

// SetSigningMethod 设定签名方式
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningSecret 设定签名方式
func SetSigningSecret(secret string) Option {
	return func(o *options) {
		o.signingSecret = []byte(secret)
	}
}

// SetExpired 设定令牌过期时长(单位秒，默认7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// SetKeyFunc 设定签名key
func SetKeyFunc(f func(context.Context, *jwt.Token, jwt.SigningMethod, interface{}) (interface{}, error)) Option {
	return func(o *options) {
		o.keyFunc = f
	}
}

// SetNewClaims 设定声明内容
func SetNewClaims(f func(context.Context, jwt.Claims, jwt.SigningMethod, interface{}) (string, error)) Option {
	return func(o *options) {
		o.signingFunc = f
	}
}

// SetTokenFunc 设定令牌Token
func SetTokenFunc(f func(context.Context) (string, error)) Option {
	return func(o *options) {
		o.tokenFunc = f
	}
}

// SetParseFunc 设定刷新者
func SetParseFunc(f func(context.Context, string) (*UserClaims, error)) Option {
	return func(o *options) {
		o.parseFunc = f
	}
}

// SetUpdateFunc 设定刷新者
func SetUpdateFunc(f func(context.Context) error) Option {
	return func(o *options) {
		o.updateFunc = f
	}
}

// SetFixClaimsFunc 设定刷新者
func SetFixClaimsFunc(f func(context.Context, *UserClaims) error) Option {
	return func(o *options) {
		o.claimsFunc = f
	}
}

//===================================================
// 分割线
//===================================================

// New 创建认证实例
func New(store store.Storer, opts ...Option) *Auther {
	o := options{
		tokenType:     "JWT",
		expired:       7200,
		signingMethod: jwt.SigningMethodHS512,
		signingFunc:   NewWithClaims,
		keyFunc:       KeyFuncCallback,
		tokenFunc:     GetBearerToken,
		updateFunc:    nil,
		parseFunc:     nil,
		claimsFunc:    nil,
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.signingSecret == nil {
		o.signingSecret = []byte(NewRandomID()) // 默认随机生成
		logger.Infof(nil, "new random signing secret: %s", o.signingSecret)
	}

	return &Auther{
		opts:  &o,
		store: store,
	}
}

// Release 释放资源
func (a *Auther) Release() error {
	return a.callStore(func(store store.Storer) error {
		return store.Close()
	})
}

//===================================================
// 分割线
//===================================================

var _ auth.Auther = &Auther{}

// Auther jwt认证
type Auther struct {
	opts  *options
	store store.Storer
}

// GetUserInfo 获取用户
func (a *Auther) GetUserInfo(c context.Context) (auth.UserInfo, error) {
	tokenString, err := a.opts.tokenFunc(c)
	if err != nil {
		return nil, err
	}

	claims := new(UserClaims)
	if a.opts.parseFunc != nil {
		claims, err = a.opts.parseFunc(c, tokenString)
	} else {
		claims, err = a.parseToken(c, tokenString)
	}
	if err != nil {
		var e *jwt.ValidationError
		if errors.As(err, &e) {
			return nil, auth.ErrInvalidToken
		}
		return nil, err
	}

	err = a.callStore(func(store store.Storer) error {
		// 反向验证该用户是否已经登出
		if exists, err := store.Check(c, claims.GetTokenID()); err != nil {
			return err
		} else if exists {
			return auth.ErrInvalidToken
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// GenerateToken 生成令牌
func (a *Auther) GenerateToken(c context.Context, user auth.UserInfo) (auth.TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.opts.expired) * time.Second).Unix()
	issuedAt := now.Unix()

	claims := NewUserInfo(user)
	claims.IssuedAt = issuedAt
	claims.NotBefore = issuedAt
	claims.ExpiresAt = expiresAt

	if a.opts.claimsFunc != nil {
		if err := a.opts.claimsFunc(c, claims); err != nil {
			return nil, err
		}
	}

	tokenString, err := a.opts.signingFunc(c, claims, a.opts.signingMethod, a.opts.signingSecret)
	if err != nil {
		return nil, err
	}

	tokenInfo := &TokenInfo{
		AccessToken: tokenString,
		TokenStatus: "ok",
		TokenType:   a.opts.tokenType,
		ExpiresAt:   expiresAt,
	}
	return tokenInfo, nil
}

// DestroyToken 销毁令牌
func (a *Auther) DestroyToken(c context.Context, user auth.UserInfo) error {
	claims, ok := user.(*UserClaims)
	if !ok {
		return auth.ErrInvalidToken
	}

	// 如果设定了存储，则将未过期的令牌放入
	return a.callStore(func(store store.Storer) error {
		expired := time.Unix(claims.ExpiresAt, 0).Sub(time.Now())
		return store.Set1(c, claims.GetTokenID(), expired)
	})
}

// UpdateAuther 更新
func (a *Auther) UpdateAuther(c context.Context) error {
	if a.opts.updateFunc != nil {
		return a.opts.updateFunc(c)
	}
	return nil
}

//===================================================
// 分割线
//===================================================

// 解析令牌
func (a *Auther) parseToken(c context.Context, tokenString string) (*UserClaims, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) { return a.keyFunc(c, t) }
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, keyFunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, auth.ErrInvalidToken
	}

	return token.Claims.(*UserClaims), nil
}

// 获取密钥
func (a *Auther) keyFunc(c context.Context, t *jwt.Token) (interface{}, error) {
	if a.opts.keyFunc == nil {
		return a.opts.signingSecret, nil
	}
	return a.opts.keyFunc(c, t, a.opts.signingMethod, a.opts.signingSecret)
}

// 调用存储方法
func (a *Auther) callStore(fn func(store.Storer) error) error {
	if store := a.store; store != nil {
		return fn(store)
	}
	return nil
}

//===================================================
// 分割线
//===================================================

// KeyFuncCallback 解析方法使用此回调函数来提供验证密钥。
// 该函数接收解析后的内容，但未验证的令牌。这使您可以在令牌的标头（例如 kid），以标识要使用的密钥。
func KeyFuncCallback(c context.Context, token *jwt.Token, method jwt.SigningMethod, secret interface{}) (interface{}, error) {
	//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//	return nil, auth.ErrInvalidToken // 无法验证
	//}
	//kid := token.Header["kid"]
	//if kid == "" {
	//	return nil, auth.ErrInvalidToken // 无法验证
	//}
	token.Method = method // 强制使用配置, 防止alg使用none而跳过验证
	return secret, nil
}

// NewWithClaims new claims
// jwt.NewWithClaims
func NewWithClaims(c context.Context, claims jwt.Claims, method jwt.SigningMethod, secret interface{}) (string, error) {
	token := &jwt.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			// "kid": "zgo123456",
		},
		Claims: claims,
		Method: method,
	}
	return token.SignedString(secret)
}

// NewRandomID new ID
func NewRandomID() string {
	// uuid, err := uuid.NewRandom()
	// if err != nil {
	// 	panic(err)
	// }
	// strid := uuid.String()
	// return strid
	return crypto.UUID(32)
}
