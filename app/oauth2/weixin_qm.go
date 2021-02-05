package oauth2

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"

	gotemplate "text/template"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/suisrc/zgo/app/model/gpa"
)

const (
	// WeixinAPI weixin api domain
	WeixinAPI = "https://api.weixin.qq.com"
)

var _ Handler = (*WeixinQm)(nil)

// WeixinQm 微信公众号
type WeixinQm struct {
	gpa.GPA              // 数据库操作
	Storer  store.Storer // 缓存器操作

	parseOnce       sync.Once            // once
	parseTemplate   *gotemplate.Template // template
	parseError      error                // error
	parseQrOnce     sync.Once            // qr once
	parseQrTemplate *gotemplate.Template // qr template
	parseQrError    error                // qr error
}

// GetUser ...
func (a *WeixinQm) GetUser(c *gin.Context, rp RequestPlatform, rt RequestToken, relation, userid string) (*UserInfo, error) {
	return nil, errors.New("no implemented")
}

// Handle handle
func (a *WeixinQm) Handle(c *gin.Context, body RequestParams, platform RequestPlatform, token1 RequestToken, oauth2 RequestOAuth2) error {
	if platform.GetAppID() == "" || platform.GetAppSecret() == "" {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-CONFIG", Other: "应用配置异常"})
	}
	// 重定向到微信服务器
	if body.GetCode() == "" {
		if strings.Contains(c.Request.UserAgent(), "MicroMessenger/") {
			return a.Connect(c, body, platform, oauth2)
		}
		return a.QrConnect(c, body, platform, oauth2)
	}
	// 通过微信服务器回调
	if body.GetState() == "" {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-STATE", Other: "状态码无效"})
	}
	if ok, err := a.Storer.Check(c, body.GetState()); err != nil || !ok {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-STATE", Other: "状态码无效"})
	}
	a.Storer.Delete(c, body.GetState()) // 删除缓存

	// 获取用户令牌（注意，这里是用户令牌， 不是应用令牌）
	token := WeixinQmAccessToken{}
	if err := token.GetAccessToken(platform.GetAppID(), platform.GetAppSecret(), body.GetCode()); err != nil {
		return err // 网络异常
	} else if token.ErrCode != 0 || token.ErrMsg != "ok" {
		return &token // 微信服务器异常
	}

	acc, err := oauth2.FindAccount("none", token.OpenID, "", "")
	if err != nil {
		return err
	}

	if acc > 0 {
		// 更新Account内容
		atoken := &AccessToken{
			Account:      acc,
			Platform:     platform.GetID(),
			AccessToken:  sql.NullString{Valid: true, String: token.AccessToken},
			ExpiresIn:    sql.NullInt64{Valid: true, Int64: int64(token.ExpiresIn)},
			ExpiresAt:    sql.NullTime{Valid: true, Time: time.Now().Add(time.Duration(token.ExpiresIn-120) * time.Second)}, // 提前2分钟失效
			RefreshToken: sql.NullString{Valid: true, String: token.RefreshToken},
			RefreshExpAt: sql.NullTime{Valid: true, Time: time.Now().Add(time.Duration(29*24) * time.Hour)}, // 提前1天失效
			Scope:        sql.NullString{Valid: true, String: token.Scope},
		}
		if err := token1.SaveAccessToken(a.Sqlx, atoken); err != nil {
			logger.Errorf(c, logger.ErrorWW(err)) // do nothing
		}
	}
	return nil
}

// Connect 应用认证
func (a *WeixinQm) Connect(c *gin.Context, body RequestParams, platform RequestPlatform, oauth2 RequestOAuth2) error {

	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(60)*time.Second)

	uri := GetRedirectURIByOAuth2Platform(c, c.Query("redirect_uri"), platform, oauth2)
	scope := body.GetScope()
	if scope == "" {
		scope = "snsapi_base"
	}
	appid := platform.GetAppID()
	// 参数
	params := helper.H{
		"appid":         appid,  // 公众号的唯一标识
		"redirect_uri":  uri,    // 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
		"response_type": "code", // 返回类型，请填写code
		"scope":         scope,  // 应用授权作用域 snsapi_base: 不弹出授权页面，直接跳转，只能获取用户openid, snsapi_userinfo: 弹出授权页面，可通过openid拿到昵称、性别、所在地。并且， 即使在未关注的情况下，只要用户授权，也能获取其信息
		"state":         state,  // 重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
	}

	a.parseOnce.Do(func() {
		// 只加载一次, 该内容是模板 解析一次即可
		url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid={{.appid}}&redirect_uri={{.redirect_uri}}&response_type=code&scope={{.scope}}&state={{.state}}#wechat_redirect"
		a.parseTemplate, a.parseError = gotemplate.New("").Parse(url)
	})

	var buf bytes.Buffer
	// if gt, err := gotemplate.New("").Parse(WeixinAuthorizeURL); err != nil {
	if gt, err := a.parseTemplate, a.parseError; err != nil {
		return err
	} else if err2 := gt.Execute(&buf, params); err2 != nil {
		return err2
	}
	localtion := buf.String()

	return &helper.ErrorRedirect{
		State:    state,
		Location: localtion,
	}
}

// QrConnect 扫码认证
// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
func (a *WeixinQm) QrConnect(c *gin.Context, body RequestParams, platform RequestPlatform, oauth2 RequestOAuth2) error {
	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(300)*time.Second) // 5分钟等待,如果5分钟没有进行扫描登陆,直接拒绝

	uri := GetRedirectURIByOAuth2Platform(c, c.Query("redirect_uri"), platform, oauth2)
	appid := platform.GetAppID()
	if result := c.Query("result"); result == "wxLogin" {
		return helper.NewSuccess(c, helper.H{
			"self_redirect": true,              // true：手机点击确认登录后可以在 iframe 内跳转到 redirect_uri，false：手机点击确认登录后可以在 top window 跳转到 redirect_uri。默认为 false。
			"id":            "login_container", // 第三方页面显示二维码的容器id
			"appid":         appid,             // 应用唯一标识，在微信开放平台提交应用审核通过后获得
			"scope":         "snsapi_login",    // 应用授权作用域，拥有多个作用域用逗号（,）分隔，网页应用目前仅填写snsapi_login即可
			"redirect_uri":  uri,               // 重定向地址，需要进行UrlEncode
			"state":         state,             // 用于保持请求和回调的状态，授权请求后原样带回给第三方。该参数可用于防止csrf攻击（跨站请求伪造攻击）
			"style":         "",                // 提供"black"、"white"可选，默认为黑色文字描述
			"href":          "",                // 自定义样式链接，第三方可根据实际需求覆盖默认样式
		})
	}
	// 参数
	params := helper.H{
		"appid":         appid,          // 公众号的唯一标识
		"redirect_uri":  uri,            // 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
		"response_type": "code",         // 返回类型，请填写code
		"scope":         "snsapi_login", // 应用授权作用域，拥有多个作用域用逗号（,）分隔，网页应用目前仅填写snsapi_login
		"state":         state,          // 重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
	}

	a.parseOnce.Do(func() {
		// 只加载一次, 该内容是模板 解析一次即可
		url := "https://open.weixin.qq.com/connect/qrconnect?appid={{.appid}}&redirect_uri={{.redirect_uri}}&response_type={{.code}}&scope={{.scope}}&state={{.state}}#wechat_redirect"
		a.parseTemplate, a.parseError = gotemplate.New("").Parse(url)
	})

	var buf bytes.Buffer
	if gt, err := a.parseTemplate, a.parseError; err != nil {
		return err
	} else if err2 := gt.Execute(&buf, params); err2 != nil {
		return err2
	}
	localtion := buf.String()

	return &helper.ErrorRedirect{
		State:    state,
		Location: localtion,
	}
}

//======================================================================
//======================================================================

// WeixinQmError {"errcode":40029,"errmsg":"invalid code"}
type WeixinQmError struct {
	ErrCode int    `json:"errocde,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (a *WeixinQmError) Error() string {
	return "[" + strconv.Itoa(a.ErrCode) + "]:" + a.ErrMsg
}

// WeixinQmAccessToken token
/*
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE"
}
*/
type WeixinQmAccessToken struct {
	WeixinQmError
	AppID        string `json:"app_id,omitempty"`        // 应用ID
	AccessToken  string `json:"access_token,omitempty"`  // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn    int    `json:"expires_in,omitempty"`    // access_token接口调用凭证超时时间，单位（秒）
	RefreshToken string `json:"refresh_token,omitempty"` // 用户刷新access_token, refresh_token有效期为30天，当refresh_token失效之后
	OpenID       string `json:"openid,omitempty"`        // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope        string `json:"scope,omitempty"`         // 用户授权的作用域，使用逗号（,）分隔
}

// GetAccessToken 获取访问令牌
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
func (a *WeixinQmAccessToken) GetAccessToken(appid, secret, code string) error {
	return gout.GET(WeixinAPI + "/sns/oauth2/access_token").
		SetQuery(gout.H{
			"appid":      appid,
			"secret":     secret,
			"code":       code,
			"grant_type": "authorization_code",
		}).
		BindJSON(a).
		Do()
}

// GetRefreshToken 获取访问令牌
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
func (a *WeixinQmAccessToken) GetRefreshToken(appid, token string) error {
	return gout.GET(WeixinAPI + "/sns/oauth2/refresh_token").
		SetQuery(gout.H{
			"appid":         appid,
			"refresh_token": token,
			"grant_type":    "refresh_token",
		}).
		BindJSON(a).
		Do()
}

// TestAccessToken 验证访问令牌是否有效
func (a *WeixinQmAccessToken) TestAccessToken() error {
	return gout.GET(WeixinAPI + "/sns/auth").
		SetQuery(gout.H{
			"access_token": a.AccessToken,
			"openid":       a.OpenID,
		}).
		BindJSON(a).
		Do()
}

// TestToken 验证令牌是否有效,会自动刷新AccessToken令牌
// 返回值代表网络异常,本身A的ErrCode反应了令牌发生的问题
func (a *WeixinQmAccessToken) TestToken() error {
	if err := a.TestAccessToken(); err != nil {
		// 确认当前令牌是否有效
		return err
	}
	if a.ErrCode == 42001 && a.AppID != "" && a.RefreshToken != "" {
		// https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Global_Return_Code.html
		// access_token 超时, 尝试刷新令牌
		if err := a.GetRefreshToken(a.AppID, a.RefreshToken); err != nil {
			return err
		}
	}
	return nil
}

// WeixinQmUserInfo 拉取用户信息(需scope为 snsapi_userinfo)
/*
{
  "openid":"OPENID",
  "nickname": "NICKNAME",
  "sex": "1",
  "province":"PROVINCE",
  "city": "CITY",
  "country":"COUNTRY",
  "headimgurl": "http://thirdwx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46",
  "privilege": ["PRIVILEGE1", "PRIVILEGE2"],
  "unionid": "o6_bmasdasdsad6_2sgVt7hMZOPfL"
}
*/
type WeixinQmUserInfo struct {
	WeixinQmError
	OpenID     string   `json:"openid,omitempty"`     // 用户的唯一标识
	Nickname   string   `json:"nickname,omitempty"`   // 用户昵称
	Sex        string   `json:"sex,omitempty"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province,omitempty"`   // 用户个人资料填写的省份
	City       string   `json:"city,omitempty"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country,omitempty"`    // 国家，如中国为CN
	HeadImgURL string   `json:"headimgurl,omitempty"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege,omitempty"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionID    string   `json:"unionid,omitempty"`    //只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段
}

// GetUserInfo 获取访问令牌
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
func (a *WeixinQmUserInfo) GetUserInfo(token, openid, lang string) error {
	if lang == "" {
		lang = "zh_CN"
	}
	return gout.GET(WeixinAPI + "/sns/userinfo").
		SetQuery(gout.H{
			"access_token": token,
			"openid":       openid,
			"lang":         lang,
		}).
		BindJSON(a).
		Do()
}
