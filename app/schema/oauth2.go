package schema

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

// OAuth2GpaPlatform 登录使用的平台
type OAuth2GpaPlatform struct {
	ID          int64          `db:"id"`           // 唯一标识
	KID         sql.NullString `db:"kid"`          // 三方标识
	Type        sql.NullString `db:"type"`         // 平台标识
	IsSign      sql.NullBool   `db:"signin"`       // 登录标识
	OrgCode     sql.NullString `db:"org_cod"`      // 组织字段
	Status      StatusType     `db:"status"`       // 状态
	AppID       sql.NullString `db:"app_id"`       // 应用标识
	AppSecret   sql.NullString `db:"app_secret"`   // 应用密钥
	AgentID     sql.NullString `db:"agent_id"`     // 代理标识
	AgentSecret sql.NullString `db:"agent_secret"` // 代理密钥
	SuiteID     sql.NullString `db:"suite_id"`     // 套件标识
	SuiteSecret sql.NullString `db:"suite_secret"` // 套件密钥
	SigninHost  sql.NullString `db:"signin_url"`   // 登录地址
	TokenKID    sql.NullString `db:"token_kid"`    // 当前令牌
	JsSecret    sql.NullString `db:"js_secret"`    // JS密钥
	StateSecret sql.NullString `db:"state_secret"` // 回调密钥
	IsCallback  sql.NullBool   `db:"callback"`     // 回调标识
	CbDomain    sql.NullString `db:"cb_domain"`    // 默认域名
	CbScheme    sql.NullString `db:"cb_scheme"`    // 默认协议
	CbEncrypt   sql.NullString `db:"cb_encrypt"`   // 加密标识
	CbToken     sql.NullString `db:"cb_token"`     // 加密令牌
	CbEncoding  sql.NullString `db:"cb_encoding"`  // 加密编码
	String1     sql.NullString `db:"string_1"`     // 备用字段
	Number1     sql.NullInt64  `db:"number_1"`     // 备用字段
	//Version      sql.NullInt64  `db:"version" set:"=version+1"`
	//CallCount    sql.NullString `db:"call_count" set:"=version+1"`
}

// QueryByKID 查询
func (a *OAuth2GpaPlatform) QueryByKID(sqlx *sqlx.DB, kid string) error {
	if kid == "" {
		return errors.New("sql: no arg in params")
	}
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}platform where kid = ?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// UpdateAndSaveByID 更新
func (a *OAuth2GpaPlatform) UpdateAndSaveByID(sqlx *sqlx.DB) error {
	tic := sqlxc.TableIdxColumn{Table: TablePrefix + "platform", IDVal: a.ID}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNilAndSet(tic, a)
	if err != nil {
		return err
	}
	_, err = sqlx.NamedExec(SQL, params)
	return err
}

// GetID ...
func (a *OAuth2GpaPlatform) GetID() int64 {
	return a.ID
}

// GetKID ...
func (a *OAuth2GpaPlatform) GetKID() int64 {
	return a.ID
}

// GetAppID ...
func (a *OAuth2GpaPlatform) GetAppID() string {
	return a.AppID.String
}

// GetAppSecret ...
func (a *OAuth2GpaPlatform) GetAppSecret() string {
	return a.AppSecret.String
}

// GetAgentID ...
func (a *OAuth2GpaPlatform) GetAgentID() string {
	return a.AgentID.String
}

// GetAgentSecret ...
func (a *OAuth2GpaPlatform) GetAgentSecret() string {
	return a.AgentSecret.String
}

// OAuth2GpaAccountToken account
type OAuth2GpaAccountToken struct {
	TokenID      string         `db:"token_kid"`
	AccountID    sql.NullInt64  `db:"account_id"`
	TokenPID     sql.NullString `db:"token_pid"`
	TokenType    sql.NullInt32  `db:"token_type"`
	Platform     sql.NullString `db:"platform_kid"`
	AccessToken  sql.NullString `db:"access_token"`
	ExpiresAt    sql.NullTime   `db:"expires_at"`
	RefreshToken sql.NullString `db:"refresh_token"`
	RefreshExpAt sql.NullTime   `db:"refresh_exp"`
	CodeToken    sql.NullString `db:"code_token"`
	CodeExpAt    sql.NullTime   `db:"code_exp"`
	CallCount    sql.NullInt64  `db:"call_count"`
	ErrCode      sql.NullString `db:"error_code"`
	ErrMessage   sql.NullString `db:"error_message"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
	String2      sql.NullString `db:"string_2"` // 扩展字段
}

// QueryByTokenKID2 kid
func (a *OAuth2GpaAccountToken) QueryByTokenKID2(sqlx *sqlx.DB, kid string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where token_kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByPlatformAndCode2 code
func (a *OAuth2GpaAccountToken) QueryByPlatformAndCode2(sqlx *sqlx.DB, platform, code string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where code_token=? and platform_kid=? order by code_exp desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, code, platform)
}

// UpdateAndSaveByTokenKID2 更新
func (a *OAuth2GpaAccountToken) UpdateAndSaveByTokenKID2(sqlx *sqlx.DB, update bool) error {
	tic := sqlxc.TableIdxColumn{Table: TablePrefix + "token", IDCol: "token_kid", IDVal: a.TokenID, Update: sql.NullBool{Valid: true, Bool: update}}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNilAndSet(tic, a)
	if err != nil {
		return err
	}
	_, err = sqlx.NamedExec(SQL, params)
	return err
}
