package schema

import "database/sql"

// OAuth2GpaPlatform 登录使用的平台
type OAuth2GpaPlatform struct {
	ID          int64          `db:"id"`           // 唯一标识
	KID         string         `db:"kid"`          // 三方标识
	Type        string         `db:"type"`         // 平台标识
	IsSign      sql.NullBool   `db:"signin"`       // 登录标识
	OrgCode     sql.NullString `db:"org_cod"`      // 组织字段
	Status      StatusType     `db:"status"`       // 状态
	AppID       sql.NullString `db:"app_id"`       // 应用标识
	AppSecret   sql.NullString `db:"app_secret"`   // 应用密钥
	AgentID     sql.NullString `db:"agent_id"`     // 代理标识
	AgentSecret sql.NullString `db:"agent_secret"` // 代理密钥
	SuiteID     sql.NullString `db:"suite_id"`     // 套件标识
	SuiteSecret sql.NullString `db:"suite_secret"` // 套件密钥
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
