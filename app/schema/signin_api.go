package schema

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password"`                    // 密码
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	KID      string `json:"kid"`                         // 授权平台
	Org      string `json:"org"`                         // 租户
	Role     string `json:"role"`                        // 角色
	Domain   string `json:"host"`                        // 域, 如果无,使用c.Reqest.Host代替
}

// SigninOfCaptcha 使用登陆发生认证信息
type SigninOfCaptcha struct {
	Mobile string `form:"mobile"` // 手机
	Email  string `form:"email"`  // 邮箱
	Openid string `form:"openid"` // openid
	KID    string `form:"kid"`    // 平台标识
	Org    string `json:"org"`    // 租户
}

// SigninOfOAuth2 登陆参数
type SigninOfOAuth2 struct {
	Code     string `form:"code"`     // 票据
	State    string `form:"state"`    // 验签
	Scope    string `form:"scope"`    // 作用域
	KID      string `form:"kid"`      // kid
	Org      string `form:"org"`      // 租户
	Role     string `form:"role"`     // 角色
	Domain   string `form:"host"`     // 域, 如果无,使用c.Reqest.Host代替
	Redirect string `form:"redirect"` // redirect
}

// SigninResult 登陆返回值
type SigninResult struct {
	TokenStatus  string        `json:"status" default:"ok"`                   // 'ok' | 'error' 不适用boolean类型是为了以后可以增加扩展
	TokenID      string        `json:"token_id,omitempty"`                    // 访问令牌ID
	AccessToken  string        `json:"access_token,omitempty"`                // 访问令牌
	TokenType    string        `json:"token_type,omitempty" default:"bearer"` // 令牌类型
	ExpiresAt    int64         `json:"expires_at,omitempty"`                  // 过期时间
	ExpiresIn    int64         `json:"expires_in,omitempty"`                  // 过期时间
	RefreshToken string        `json:"refresh_token,omitempty"`               // 刷新令牌
	RefreshExpAt int64         `json:"refresh_expires,omitempty"`             // 刷新令牌过期时间
	Redirect     string        `json:"redirect_uri,omitempty"`                // redirect_uri
	Message      string        `json:"message,omitempty"`                     // 消息,有限显示 // Message 和 Datas 一般用户发生异常后回显
	Params       []interface{} `json:"params,omitempty"`                      // 多租户多角色的时候，返回角色，重新确认登录
}
