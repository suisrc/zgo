package schema

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

// ClientGpaWebToken 登录使用的平台
type ClientGpaWebToken struct {
	KID          string         `db:"kid"`          // 三方标识
	OrgCode      sql.NullString `db:"org_cod"`      // 平台标识
	Target       sql.NullInt64  `db:"target"`       // 终端标识
	Type         sql.NullString `db:"type"`         // 终端类型
	Status       StatusType     `db:"status"`       // 状态
	JwtExpired   sql.NullInt64  `db:"jwt_expired"`  // 令牌有效期
	JwtRefresh   sql.NullInt64  `db:"jwt_refresh"`  // 令牌有效期
	JwtType      sql.NullString `db:"jwt_type"`     // 令牌类型
	JwtMethod    sql.NullString `db:"jwt_method"`   // 令牌方法
	JwtSecret    sql.NullString `db:"jwt_secret"`   // 令牌密钥
	JwtGetter    sql.NullString `db:"jwt_getter"`   // 令牌获取方法
	JwtIssuer    sql.NullString `db:"jwt_issuer"`   // 令牌签发平台
	JwtAudience  sql.NullString `db:"jwt_audience"` // 令牌接受平台
	JwtSigninURL sql.NullString `db:"signin_url"`   // 登陆地址
	JwtSigninCHK sql.NullString `db:"signin_check"` // 登陆确认
	String1      sql.NullInt64  `db:"string_1"`     // 备用字段
	Number1      sql.NullInt64  `db:"number_1"`     // 备用字段
	CreatedAt    sql.NullTime   `db:"created_at"`   // 创建时间
	UpdatedAt    sql.NullTime   `db:"updated_at"`   // 更新时间
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
	SecretByte   []byte         `db:"-"`
}

// QueryByOrg 查询
func (a *ClientGpaWebToken) QueryByOrg(sqlx *sqlx.DB, org string) error {
	if org == "" {
		return errors.New("sql: no arg in params")
	}
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}web_token where org_cod = ? and type = 'org' and status = 1 order by created_at desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, org)
}

// QueryByUsr 查询
func (a *ClientGpaWebToken) QueryByUsr(sqlx *sqlx.DB, tid int64) error {
	if tid <= 0 {
		return errors.New("sql: no arg in params")
	}
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}web_token where target = ? and type = 'usr' and status = 1 order by created_at desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, tid)
}

// QueryByApp 查询
func (a *ClientGpaWebToken) QueryByApp(sqlx *sqlx.DB, tid int64) error {
	if tid <= 0 {
		return errors.New("sql: no arg in params")
	}
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}web_token where target = ? and type = 'app' and status = 1 order by created_at desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, tid)
}

// QueryByKID 查询
func (a *ClientGpaWebToken) QueryByKID(sqlx *sqlx.DB, kid string) error {
	if kid == "" {
		return errors.New("sql: no arg in params")
	}
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}web_token where kid = ?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}
