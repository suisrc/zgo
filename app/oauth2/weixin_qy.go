package oauth2

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"sync"
	gotemplate "text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/modules/crypto"
	"github.com/suisrc/zgo/modules/helper"
	"github.com/suisrc/zgo/modules/store"
)

const (
	// WeixinQyAPI weixin api domain
	WeixinQyAPI = "https://qyapi.weixin.qq.com"
)

var _ Handler = &WeixinQy{}

// WeixinQy 微信企业号
type WeixinQy struct {
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
func (a *WeixinQy) GetUser(c *gin.Context, rp RequestPlatform, rt RequestToken, relation, userid string) (*UserInfo, error) {
	return nil, errors.New("no implemented")
}

// Handle handle
func (a *WeixinQy) Handle(c *gin.Context, body RequestParams, platform RequestPlatform, token1 RequestToken, oauth2 RequestOAuth2) error {
	if platform.GetAppID() == "" || platform.GetAgentID() == "" || (platform.GetAppSecret() == "" && platform.GetAgentSecret() == "") {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-CONFIG", Other: "应用配置异常"})
	}
	// 重定向选择
	if body.GetCode() == "" {
		if strings.Contains(c.Request.UserAgent(), "MicroMessenger/") {
			return a.Connect(c, body, platform, oauth2)
		}
		return a.QrConnect(c, body, platform, oauth2)
	}
	// 微信服务器重定向
	if body.GetState() == "" {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-STATE", Other: "状态码无效"})
	}
	if ok, err := a.Storer.Check(c, body.GetState()); err != nil || !ok {
		return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-STATE", Other: "状态码无效"})
	}
	a.Storer.Delete(c, body.GetState()) // 删除缓存

	// 获取当前用户信息
	user := WeixinQyUserInfo{}
	if err := WeixinQyExecWithAccessToken(c, a.GPA, a.Storer, platform, token1, func(token string) error {
		if err := user.GetUserInfo(token, body.GetCode()); err != nil {
			return err
		} else if user.ErrCode != 0 || user.ErrMsg != "ok" {
			return &user // 微信服务器异常, 当发生42001异常,会直接获取令牌重试一次
		} else if user.OpenID == "" && user.UserID != "" {
			// 成员用户， 强制取成员openid， 归一操作
			openid := WeixinQyOpenid{}
			if err := openid.ConvertToOpenid(token, user.UserID); err == nil && openid.ErrCode == 0 {
				user.OpenID = openid.OpenID
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if user.UserID != "" {
		_, err := oauth2.FindAccount("member", user.OpenID, user.UserID, user.DeviceID)
		return err
	}
	if user.OpenID != "" {
		// return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-MEMBER", Other: "非成员用户"})
		if user.ExternalUserid != "" {
			// 外部联系人
			_, err := oauth2.FindAccount("external", user.OpenID, user.ExternalUserid, user.DeviceID)
			return err
		}
		_, err := oauth2.FindAccount("none", user.OpenID, "", user.DeviceID)
		return err
	}
	return helper.New0Error(c, helper.ShowWarn, &i18n.Message{ID: "WARN-OAUTH2-OPENID", Other: "无法获取用户的OPEN-ID"})
}

// Connect 应用认证
func (a *WeixinQy) Connect(c *gin.Context, body RequestParams, platform RequestPlatform, oauth2 RequestOAuth2) error {
	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(60)*time.Second)

	uri := GetRedirectURIByOAuth2Platform(c, c.Query("redirect_uri"), platform, oauth2)
	scope := "snsapi_base"
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
func (a *WeixinQy) QrConnect(c *gin.Context, body RequestParams, platform RequestPlatform, oauth2 RequestOAuth2) error {
	// 未取得code内容,需要重定向回到微信服务器
	state := crypto.UUID(32)
	a.Storer.Set1(c, state, time.Duration(300)*time.Second) // 5分钟等待,如果5分钟没有进行扫描登陆,直接拒绝

	uri := GetRedirectURIByOAuth2Platform(c, c.Query("redirect_uri"), platform, oauth2)
	appid := platform.GetAppID()
	agentid := platform.GetAgentID()
	// 参数
	params := helper.H{
		"appid":        appid,   // 公众号的唯一标识
		"agentid":      agentid, // 返回类型，请填写code
		"redirect_uri": uri,     // 授权后重定向的回调链接地址， 请使用 urlEncode 对链接进行处理
		"state":        state,   // 重定向后会带上state参数，开发者可以填写a-zA-Z0-9的参数值，最多128字节
	}

	a.parseOnce.Do(func() {
		// 只加载一次, 该内容是模板 解析一次即可
		url := "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid={{.appid}}&agentid={{.agentid}}&redirect_uri={{.redirect_uri}}&state={{.state}}"
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

//===================================================================================================AccessToken-START

// WeixinQyExecWithAccessToken 执行访问, 带有token
func WeixinQyExecWithAccessToken(c context.Context, GPA gpa.GPA, Storer store.Storer, platform RequestPlatform, token1 RequestToken, fn func(token string) error) error {
	GetAccessToken := func() (*AccessToken, error) {
		token := WeixinQyAccessToken{}
		// 优先使用代理令牌， 如果代理令牌不存在， 使用应用令牌（理论上， 应用令牌具有更高的权限）
		secret := helper.IfString(platform.GetAgentSecret() != "", platform.GetAgentSecret(), platform.GetAppSecret())
		if err := token.GetAccessToken(platform.GetAppID(), secret); err != nil {
			return nil, err // 网络异常
		} else if token.ErrCode != 0 {
			return nil, &token // 微信服务器异常
		} else if token.AccessToken == "" {
			return nil, errors.New(token.ErrMsg) // 未知异常
		}
		tid, tok := helper.GetCtxValueToString(c, helper.ResTokenKey)
		return &AccessToken{
			Platform:    platform.GetID(),
			TokenID:     sql.NullString{Valid: tok, String: tid},
			AccessToken: sql.NullString{Valid: true, String: token.AccessToken},
			ExpiresIn:   sql.NullInt64{Valid: true, Int64: int64(token.ExpiresIn)},
			ExpiresAt:   sql.NullTime{Valid: true, Time: time.Now().Add(time.Duration(token.ExpiresIn-120) * time.Second)}, // 有效期缩短120秒
		}, nil
	}
	manager := &TokenManager{
		GPA:          GPA,
		TokenKey:     "token:weixin:qy:platform_id_" + strconv.Itoa(int(platform.GetID())),
		Platform:     platform.GetID(),
		Storer:       Storer,
		OAuth2Handle: token1,
		NewTokenFunc: GetAccessToken,
		MaxCacheIdle: 300,
		MinCacheIdle: 60,
		NewCacheIdle: 1800,
	}
	// 执行内容
	execFunc := func(token string) (bool, interface{}, error) {
		if err := fn(token); err != nil {
			if token, ok := err.(*WeixinQyAccessToken); ok {
				if token.ErrCode == 42001 {
					// 	access_token有时效性，需要重新获取一次
					return true, nil, nil
				}
			}
			return false, nil, err
		}
		return false, nil, nil
	}
	// 执行
	_, err := ExecWithAccessTokenX(c, execFunc, manager)
	return err
}

//===================================================================================================AccessToken-END

// WeixinQyError {"errcode":40029,"errmsg":"invalid code"}
type WeixinQyError struct {
	ErrCode int    `json:"errocde,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (a *WeixinQyError) Error() string {
	return "[" + strconv.Itoa(a.ErrCode) + "]:" + a.ErrMsg
}

// WeixinQyAccessToken token
/*
{
   "errcode": 0,
   "errmsg": "ok",
   "access_token": "accesstoken000001",
   "expires_in": 7200
}
*/
type WeixinQyAccessToken struct {
	WeixinQyError
	AccessToken string `json:"access_token,omitempty"` // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn   int    `json:"expires_in,omitempty"`   // access_token接口调用凭证超时时间，单位（秒）
}

// GetAccessToken 获取访问令牌
// https://work.weixin.qq.com/api/doc/90000/90135/91039
func (a *WeixinQyAccessToken) GetAccessToken(appid, secret string) error {
	return gout.GET(WeixinQyAPI + "/cgi-bin/gettoken").
		SetQuery(gout.H{
			"corpid":     appid,
			"corpsecret": secret,
		}).
		BindJSON(a).
		Do()
}

// WeixinQyUserInfo user info
/*
a) 当用户为企业成员时返回示例如下：
{
   "errcode": 0,
   "errmsg": "ok",
   "UserId":"USERID",
   "DeviceId":"DEVICEID"
}
b) 非企业成员授权时返回示例如下：
{
   "errcode": 0,
   "errmsg": "ok",
   "OpenId":"OPENID",
   "DeviceId":"DEVICEID",
   "external_userid":"EXTERNAL_USERID"
}
*/
type WeixinQyUserInfo struct {
	WeixinQyError
	OpenID         string `json:"OpenId,omitempty"`          // 非企业成员的标识，对当前企业唯一。不超过64字节
	UserID         string `json:"UserId,omitempty"`          // 成员UserID。若需要获得用户详情信息，可调用通讯录接口：读取成员。如果是互联企业，则返回的UserId格式如：CorpId/userid
	DeviceID       string `json:"DeviceId,omitempty"`        // 手机设备号(由企业微信在安装时随机生成，删除重装会改变，升级不受影响)
	ExternalUserid string `json:"external_userid,omitempty"` // 外部联系人id，当且仅当用户是企业的客户，且跟进人在应用的可见范围内时返回。
}

// GetUserInfo 获取访问令牌
// https://work.weixin.qq.com/api/doc/90000/90135/91023
// code : 通过成员授权获取到的code，最大为512字节。每次成员授权带上的code将不一样，code只能使用一次，5分钟未被使用自动过期。
/*
权限说明：
跳转的域名须完全匹配access_token对应应用的可信域名，否则会返回50001错误。
*/
func (a *WeixinQyUserInfo) GetUserInfo(token, code string) error {
	return gout.GET(WeixinQyAPI + "/cgi-bin/user/getuserinfo").
		SetQuery(gout.H{
			"access_token": token,
			"code":         code,
		}).
		BindJSON(a).
		Do()
}

// WeixinQyOpenid ...
type WeixinQyOpenid struct {
	WeixinQyError
	OpenID string `json:"openid,omitempty"` // 非企业成员的标识，对当前企业唯一。不超过64字节
}

// ConvertToOpenid 获取访问令牌
// https://work.weixin.qq.com/api/doc/90000/90135/90202
/*
权限说明：
成员必须处于应用的可见范围内

该接口使用场景为企业支付，在使用企业红包和向员工付款时，需要自行将企业微信的userid转成openid。
注：需要成员使用微信登录企业微信或者关注微工作台（原企业号）才能转成openid;
如果是外部联系人，请使用外部联系人openid转换转换openid
*/
func (a *WeixinQyOpenid) ConvertToOpenid(token, userid string) error {
	return gout.POST(WeixinQyAPI + " https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid").
		SetQuery(gout.H{
			"access_token": token,
		}).
		SetBody(gout.H{
			"userid": userid,
		}).
		BindJSON(a).
		Do()
}
