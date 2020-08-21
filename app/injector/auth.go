package injector

import (
	"context"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/auth/jwt"
	"github.com/suisrc/zgo/modules/auth/jwt/store/buntdb"
)

// NewAuther of auth.Auther
// 授权认证使用的auther内容
func NewAuther(gpa service.GPA) auth.Auther {
	store, err := buntdb.NewStore(":memory:") // 使用内存缓存
	if err != nil {
		panic(err)
	}
	//  secret := config.C.JWTAuth.SigningSecret
	//  if secret == "" {
	//  	secret = auth.UUID(128)
	//  	logger.Infof(nil, "jwt secret: %s", secret)
	//  }

	//  auther := jwt.New(store,
	//  	jwt.SetSigningSecret(secret), // 注册令牌签名密钥
	//  )

	auther := jwt.New(store,
		jwt.SetKeyFunc(KeyFuncCallback),
		jwt.SetNewClaims(NewWithClaims),
	)

	return auther
}

// KeyFuncCallback 解析方法使用此回调函数来提供验证密钥。
// 该函数接收解析后的内容，但未验证的令牌。这使您可以在令牌的标头（例如 kid），以标识要使用的密钥。
func KeyFuncCallback(token *jwtgo.Token, method jwtgo.SigningMethod, secret interface{}) (interface{}, error) {
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
func NewWithClaims(c context.Context, claims jwtgo.Claims, method jwtgo.SigningMethod) (*jwtgo.Token, error) {
	return &jwtgo.Token{
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": method.Alg(),
			// "kid": "zgo123456",
		},
		Claims: claims,
		Method: method,
	}, nil
}
