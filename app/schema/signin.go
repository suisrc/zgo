package schema

import (
	"database/sql"
	"strings"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/modules/auth"
)

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password"`                    // 密码
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	KID      string `json:"kid"`                         // 授权平台
	Role     string `json:"role"`                        // 角色
	Client   string `json:"client"`                      // 子应用ID
	Domain   string `json:"host"`                        // 域, 如果无,使用c.Reqest.Host代替
}

// SigninOfCaptcha 使用登陆发生认证信息
type SigninOfCaptcha struct {
	Mobile string `form:"mobile"` // 手机
	Email  string `form:"email"`  // 邮箱
	Openid string `form:"openid"` // openid
	KID    string `form:"kid"`    // 平台标识
}

// SigninOfOAuth2 登陆参数
type SigninOfOAuth2 struct {
	Code     string `form:"code"`         // 票据
	State    string `form:"state"`        // 验签
	Scope    string `form:"scope"`        // 作用域
	KID      string `form:"kid"`          // kid
	Role     string `form:"role"`         // 角色
	Client   string `form:"client"`       // 子应用ID
	Domain   string `form:"host"`         // 域, 如果无,使用c.Reqest.Host代替
	Redirect string `form:"redirect_uri"` // redirect_uri
}

// SigninResult 登陆返回值
type SigninResult struct {
	TokenStatus  string `json:"status" default:"ok"`                   // 'ok' | 'error' 不适用boolean类型是为了以后可以增加扩展
	AccessToken  string `json:"access_token,omitempty"`                // 访问令牌
	TokenType    string `json:"token_type,omitempty" default:"bearer"` // 令牌类型
	ExpiresAt    int64  `json:"expires_at,omitempty"`                  // 过期时间
	ExpiresIn    int64  `json:"expires_in,omitempty"`                  // 过期时间
	RefreshToken string `json:"refresh_token,omitempty"`               // 刷新令牌
	Redirect     string `json:"redirect_uri,omitempty"`                // redirect_uri
	// Message 和 Roles 一般用户发生异常后回显
	Message string        `json:"message,omitempty"` // 消息,有限显示
	Roles   []interface{} `json:"roles,omitempty"`   // 多角色的时候，返回角色，重新确认登录
}

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	UserName  string
	UserID    string
	RoleID    string
	TokenID   string
	Issuer    string
	Audience  string
	AccountID string
}

// GetUserName 用户名
func (s *SigninUser) GetUserName() string {
	return s.UserName
}

// GetUserID 用户ID
func (s *SigninUser) GetUserID() string {
	return s.UserID
}

// GetRoleID 角色ID
func (s *SigninUser) GetRoleID() string {
	return s.RoleID
}

// SetRoleID 角色ID
func (s *SigninUser) SetRoleID(nrole string) string {
	orole := s.RoleID
	s.RoleID = nrole
	return orole
}

// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
func (s *SigninUser) GetTokenID() string {
	return s.TokenID
}

// GetAccountID token
func (s *SigninUser) GetAccountID() string {
	return s.AccountID
}

// GetProps 获取私有属性,该内容会被加密, 注意:内容敏感,不要存储太多的内容
func (s *SigninUser) GetProps() (interface{}, bool) {
	return nil, false
}

// GetIssuer 令牌签发者
func (s *SigninUser) GetIssuer() string {
	return s.Issuer
}

// GetAudience 令牌接收者
func (s *SigninUser) GetAudience() string {
	return s.Audience
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaUser user
type SigninGpaUser struct {
	ID     int    `db:"id" json:"-"`
	KID    string `db:"kid" json:"id"`
	Name   string `db:"name" json:"name"`
	Status bool   `db:"status" json:"-"`
}

// QueryByID sql select
func (a *SigninGpaUser) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name, status from {{TP}}user where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// SigninGpaRole role
type SigninGpaRole struct {
	ID     int     `db:"id" json:"-"`
	KID    string  `db:"kid" json:"id"`
	Name   string  `db:"name" json:"name"`
	Domain *string `db:"domain" json:"domain"`
	//Domain sql.NullString `db:"domain" json:"domain"`
}

// QueryByID sql select
func (a *SigninGpaRole) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name, domain from {{TP}}role where id=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByKID sql select
func (a *SigninGpaRole) QueryByKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select id, kid, name, domain from {{TP}}role where kid=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByUserID sql select
func (a *SigninGpaRole) QueryByUserID(sqlx *sqlx.DB, dest *[]SigninGpaRole, userid int) error {
	SQL := "select r.id, r.kid, r.name, r.domain from {{TP}}user_role ur inner join {{TP}}role r on r.id=ur.role_id where ur.user_id=? and r.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL, userid)
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaAccount account
type SigninGpaAccount struct {
	ID           int            `db:"id"`
	PID          sql.NullInt64  `db:"pid"`
	Account      string         `db:"account"`
	AccountType  int            `db:"account_typ"`
	AccountKind  sql.NullString `db:"account_kid"`
	Password     sql.NullString `db:"password"`
	PasswordSalt sql.NullString `db:"password_salt"`
	PasswordType sql.NullString `db:"password_type"`
	VerifySecret sql.NullString `db:"verify_secret"`
	UserID       int            `db:"user_id"`
	RoleID       sql.NullInt64  `db:"role_id"`

	// SQLX1 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
	// SQLX2 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
}

// QueryByID 查询
func (a *SigninGpaAccount) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}account where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByAccount sql select
func (a *SigninGpaAccount) QueryByAccount(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("select " + sqlxc.SelectColumns(a, ""))
	sqr.WriteString(" from {{TP}}account")
	sqr.WriteString(" where account=? and account_typ=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		sqr.WriteString(" and account_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and account_kid is null")
	}
	sqr.WriteString(" and status=1")
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// UpdateVerifySecret update verify secret
func (a *SigninGpaAccount) UpdateVerifySecret(sqlx *sqlx.DB) error {
	SQL := "update {{TP}}account set verify_secret=? where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	_, err := sqlx.Exec(SQL, a.VerifySecret.String, a.ID)
	return err
}

// SigninGpaAccountOA2 account
type SigninGpaAccountOA2 struct {
	ID           int            `db:"id"`
	AccessToken  sql.NullString `db:"oa2_token"`   // oauth2令牌
	ExpiresAt    sql.NullTime   `db:"oa2_expired"` // oauth2过期时间
	RefreshToken sql.NullString `db:"oa2_refresh"` // 刷新令牌
	Scope        sql.NullString `db:"oa2_scope"`   // 授权作用域
}

// UpdateOAuth2Info update
func (a *SigninGpaAccountOA2) UpdateOAuth2Info(sqlx *sqlx.DB) error {
	IDC := sqlxc.IDC{ID: int64(a.ID)}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNil(TablePrefix+"account", "id", IDC, a)
	if err != nil {
		return err
	}

	res, err := sqlx.NamedExec(SQL, params)
	if err != nil {
		return err
	}
	if IDC.ID == 0 {
		ID, _ := res.LastInsertId()
		a.ID = int(ID)
	}
	return nil
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaAccountToken account
type SigninGpaAccountToken struct {
	ID           int            `db:"id"`
	AccountID    int            `db:"account_id"`
	UserKID      string         `db:"user_kid"`
	TokenID      string         `db:"token_kid"`
	ClientID     sql.NullInt64  `db:"client_id"`
	ClientKID    sql.NullString `db:"client_kid"`
	RoleKID      sql.NullString `db:"role_kid"`
	LastIP       sql.NullString `db:"last_ip"`
	LastAt       sql.NullTime   `db:"last_at"`
	LimitExp     sql.NullTime   `db:"limit_exp"`
	LimitKey     sql.NullString `db:"limit_key"`
	Mode         sql.NullString `db:"mode"`
	ExpiresAt    sql.NullInt64  `db:"expires_at"`
	AccessToken  sql.NullString `db:"access_token"`
	RefreshToken sql.NullString `db:"refresh_token"`
	RefreshCount sql.NullInt64  `db:"refresh_count" set:"=refresh_count+1"`
	Status       sql.NullBool   `db:"status"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
}

// QueryByAccountAndClient kid
func (a *SigninGpaAccountToken) QueryByAccountAndClient(sqlx *sqlx.DB, accountID, clientID int, clientIP string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}account_token where account_id=?"
	params := []interface{}{accountID}
	if clientID > 0 {
		SQL += " and client_id=?"
		params = append(params, clientID)
	} else {
		SQL += " and client_id is null"
	}
	if clientIP != "" {
		SQL += " and last_ip=?"
		params = append(params, clientIP)
	}
	SQL += " order by expires_at desc limit 1"

	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// QueryByAccountAndClientK kid
func (a *SigninGpaAccountToken) QueryByAccountAndClientK(sqlx *sqlx.DB, accountID int, clientKID, clientIP string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}account_token where account_id=?"
	params := []interface{}{accountID}
	if clientKID == "" {
		SQL += " and client_kid is null"
	} else {
		SQL += " and client_kid=?"
		params = append(params, clientKID)
	}
	if clientIP != "" {
		SQL += " and last_ip=?"
		params = append(params, clientIP)
	}
	SQL += " order by expires_at desc limit 1"

	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// QueryByRefreshToken tid
func (a *SigninGpaAccountToken) QueryByRefreshToken(sqlx *sqlx.DB, token string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}account_token where refresh_token=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, token)
}

// UpdateAndSaveByAccountAndClient 更新
func (a *SigninGpaAccountToken) UpdateAndSaveByAccountAndClient(sqlx *sqlx.DB) (int64, error) {
	IDC := sqlxc.IDC{}
	if a.TokenID != "" {
		SQL := "select id from {{TP}}account_token where token_kid=?"
		SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
		sqlx.Get(&IDC, SQL, a.TokenID)
	}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNilAndSet(TablePrefix+"account_token", "id", IDC, a)
	if err != nil {
		return 0, err
	}
	// tx := sqlx.MustBegin()
	// tx.MustExec(SQL, params)
	// tx.Commit()
	res, err := sqlx.NamedExec(SQL, params)
	if err != nil {
		return 0, err
	}
	if IDC.ID > 0 {
		return IDC.ID, nil
	}
	return res.LastInsertId()
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaOAuth2Platfrm 第三方登陆实体
type SigninGpaOAuth2Platfrm struct {
	ID           int            `db:"id"`
	KID          string         `db:"kid"`           // 三方标识
	Platform     string         `db:"platform"`      // 平台标识, 主要用户识别操作句柄
	AppID        sql.NullString `db:"app_id"`        // 应用标识
	AppSecret    sql.NullString `db:"app_secret"`    // 应用密钥
	Avatar       sql.NullString `db:"avatar"`        // 平台头像
	Description  sql.NullString `db:"description"`   // 平台描述
	Status       sql.NullBool   `db:"status"`        // 状态
	Signin       sql.NullBool   `db:"signin"`        // 可登陆
	AgentID      sql.NullString `db:"agent_id"`      // 代理商标识
	AgentSecret  sql.NullString `db:"agent_secret"`  // 代理商密钥
	SuiteID      sql.NullString `db:"suite_id"`      // 套件标识
	SuiteSecret  sql.NullString `db:"suite_secret"`  // 套件密钥
	AuthorizeURL sql.NullString `db:"authorize_url"` // 认证地址
	TokenURL     sql.NullString `db:"token_url"`     // 令牌地址
	ProfileURL   sql.NullString `db:"profile_url"`   // 个人资料地址
	SigninURL    sql.NullString `db:"signin_url"`    // 上游应用无法获取https时候替代方案
	JsSecret     sql.NullString `db:"js_secret"`     // javascript密钥
	StateSecret  sql.NullString `db:"state_secret"`  // 回调state密钥
	Callback     sql.NullBool   `db:"callback"`      // 是否支持回调
	CbDomain     sql.NullString `db:"cb_domain"`     // 默认域名
	CbScheme     sql.NullString `db:"cb_scheme"`     // 默认协议
	CbEncrypt    sql.NullString `db:"cb_encrypt"`    // 回调是否加密
	CbToken      sql.NullString `db:"cb_token"`      // 加密令牌
	CbEncoding   sql.NullString `db:"cb_encoding"`   // 加密编码
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
}

// QueryByID id
func (a *SigninGpaOAuth2Platfrm) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_platform where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByKID kid
func (a *SigninGpaOAuth2Platfrm) QueryByKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_platform where kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}
