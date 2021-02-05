package schema

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password"`                    // 密码
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	Scope    string `json:"scope"`                       // 作用域
	Platform string `json:"p"`                           // 授权平台
	OrgCode  string `json:"g"`                           // 租户
	WebToken string `json:"w"`                           // JWT令牌密钥， 高级用法， 使用非系统默认JWT令牌
}

// SigninOfCaptcha 使用登陆发生认证信息
type SigninOfCaptcha struct {
	Mobile   string `form:"mobile"` // 手机
	Email    string `form:"email"`  // 邮箱
	Openid   string `form:"openid"` // openid
	Platform string `form:"p"`      // 授权平台
	OrgCode  string `form:"g"`      // 租户
}

// SigninOfOAuth2 登陆参数
type SigninOfOAuth2 struct {
	Code     string `form:"code"`         // 票据
	State    string `form:"state"`        // 验签
	Scope    string `form:"scope"`        // 作用域
	Platform string `form:"p" uri:"p"`    // 授权平台
	OrgCode  string `form:"g" uri:"g"`    // 租户
	WebToken string `form:"w" uri:"w"`    // JWT令牌密钥， 高级用法， 使用非系统默认JWT令牌
	Redirect string `form:"redirect_uri"` // 重定向地址
}

// GetCode ...
func (a *SigninOfOAuth2) GetCode() string {
	return a.Code
}

// GetState ...
func (a *SigninOfOAuth2) GetState() string {
	return a.Scope
}

// GetScope ...
func (a *SigninOfOAuth2) GetScope() string {
	return a.Scope
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
