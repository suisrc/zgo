package oauth2

import (
	"bytes"
	"database/sql"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store"

	gotemplate "text/template"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
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

// Handle handle
func (a *WeixinQm) Handle(c *gin.Context, b *schema.SigninOfOAuth2, o2p *schema.SigninGpaOAuth2Platfrm, acc *schema.SigninGpaAccount) error {
	// 验证配置
	if !o2p.AppID.Valid || !o2p.AppSecret.Valid {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-CONFIG", Other: "应用配置异常"})
	}
	if !o2p.Status.Valid || !o2p.Status.Bool {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-DISABLE", Other: "应用被禁用"})
	}
	if !o2p.Signin.Valid || !o2p.Signin.Bool {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NOSIGNIN", Other: "应用无法授权"})
	}
	// 重定向到微信服务器
	if b.Code == "" {
		if strings.Contains(c.Request.UserAgent(), "MicroMessenger/") {
			return a.Connect(c, b, o2p)
		}
		return a.QrConnect(c, b, o2p)
	}
	// 通过微信服务器回调
	if b.State == "" {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-STATE", Other: "状态码无效"})
	}
	if err := a.Storer.Delete(c, b.State); err != nil {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-STATE", Other: "状态码无效"})
	}
	// o2d := WeixinOAuth2Data{}
	// if data, ok, err := a.Storer.Get(c, b.State); err != nil {
	// 	return err
	// } else if !ok {
	// 	return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-STATE", Other: "状态码无效"})
	// } else if err2 := helper.JSONUnmarshal([]byte(data), &o2d); err2 != nil {
	// 	return err2
	// }
	// b.Role = o2d.Role
	// b.Client = o2d.Client
	// b.Domain = o2d.Domain
	// b.Redirect = o2d.Redirect

	// 获取用户令牌（注意，这里是用户令牌， 不是应用令牌）
	token := WeixinQmAccessToken{}
	if err := token.GetAccessToken(o2p.AppID.String, o2p.AppSecret.String, b.Code); err != nil {
		return err // 网络异常
	} else if token.ErrCode != 0 || token.ErrMsg != "ok" {
		return &token // 微信服务器异常
	}
	// 查询当前用户
	if err := acc.QueryByAccount(a.Sqlx, token.OpenID, int(schema.ATOpenid), o2p.KID); err != nil {
		if sqlxc.IsNotFound(err) {
			return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NOBIND", Other: "用户未绑定"})
		}
		return err
	}
	if acc.ID == 0 {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-SIGNIN-OAUTH2-NOBIND", Other: "用户未绑定"})
	}
	// 更新Account内容
	oa2 := &schema.SigninGpaAccountOA2{
		ID:           acc.ID,
		AccessToken:  sql.NullString{Valid: true, String: token.AccessToken},
		ExpiresAt:    sql.NullTime{Valid: true, Time: time.Now().Add(time.Duration(token.ExpiresIn-30) * time.Second)},
		RefreshToken: sql.NullString{Valid: true, String: token.RefreshToken},
		Scope:        sql.NullString{Valid: true, String: token.Scope},
	}
	if err := oa2.UpdateOAuth2Info(a.Sqlx); err != nil {
		// do nothing
		logger.Errorf(c, logger.ErrorWW(err))
	}

	return nil
}

// Connect 应用认证
func (a *WeixinQm) Connect(c *gin.Context, b *schema.SigninOfOAuth2, o2p *schema.SigninGpaOAuth2Platfrm) error {

	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(60)*time.Second)
	// o2d := WeixinOAuth2Data{
	// 	SigninOfOAuth2: *b,
	// }
	// if data, err := helper.JSONMarshal(&o2d); err == nil {
	// 	a.Storer.Set(c, state, string(data), time.Duration(30)*time.Second)
	// } else {
	// 	return err
	// }
	uri := GetRedirectURIByOAuth2Platfrm(c, o2p)
	uri = url.QueryEscape(uri) // 进行URL编码
	scope := b.Scope
	if scope == "" {
		scope = "snsapi_base"
	}
	appid := o2p.AppID.String
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
		Location: localtion,
	}
}

// QrConnect 扫码认证
// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
func (a *WeixinQm) QrConnect(c *gin.Context, b *schema.SigninOfOAuth2, o2p *schema.SigninGpaOAuth2Platfrm) error {
	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(300)*time.Second) // 5分钟等待,如果5分钟没有进行扫描登陆,直接拒绝

	uri := GetRedirectURIByOAuth2Platfrm(c, o2p)
	uri = url.QueryEscape(uri) // 进行URL编码
	appid := o2p.AppID.String
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
		Location: localtion,
	}
}

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
	RefreshToken string `json:"refresh_token,omitempty"` // 用户刷新access_token
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
