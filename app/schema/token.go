package schema

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

// AccountOAuth2Token account
type AccountOAuth2Token struct {
	TokenID      string         `db:"token_kid"`     // 令牌标识
	AccessToken  sql.NullString `db:"access_token"`  // oauth2令牌
	ExpiresAt    sql.NullTime   `db:"expires_at"`    // oauth2过期时间
	ExpiresIn    sql.NullInt64  `db:"expires_in"`    // 有限期间隔
	RefreshToken sql.NullString `db:"refresh_token"` // 刷新令牌
	Scope        sql.NullString `db:"token_scope"`   // 授权作用域
}

// UpdateToken update
func (a *AccountOAuth2Token) UpdateToken(sqlx *sqlx.DB) error {
	tic := sqlxc.TableIdxColumn{Table: TablePrefix + "token_oauth2", IDCol: "token_kid", IDVal: a.TokenID}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNil(tic, a)
	if err != nil {
		return err
	}
	_, err = sqlx.NamedExec(SQL, params)
	return err
}

// ServerOAuth2Token 第三方登陆实体
type ServerOAuth2Token struct {
	TokenID      string         `db:"token_kid"`     // 令牌标识
	PlatformID   int            `db:"oauth2_id"`     // 平台
	AccessToken  sql.NullString `db:"access_token"`  // 访问令牌
	ExpiresAt    sql.NullTime   `db:"expires_at"`    // 凭据过期时间
	ExpiresIn    sql.NullInt64  `db:"expires_in"`    // 有限期间隔
	RefreshToken sql.NullString `db:"refresh_token"` // 刷新令牌
	RefreshCount sql.NullString `db:"refresh_count"` // 刷新次数
	SyncLock     sql.NullTime   `db:"sync_lock"`     // 同步锁
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	//TokenType    sql.NullString `db:"token_type"`
	//Version      sql.NullInt64  `db:"version" set:"=version+1"`
	//CallCount    sql.NullString `db:"call_count" set:"=version+1"`
}

// QueryByTokenID kid
func (a *ServerOAuth2Token) QueryByTokenID(sqlx *sqlx.DB, id string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token_oauth2 where token_kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByPlatform platform
func (a *ServerOAuth2Token) QueryByPlatform(sqlx *sqlx.DB, platform int) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token_oauth2 where oauth2_id=? order by expires_at desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, platform)
}

// QueryByPlatformMust platform
func (a *ServerOAuth2Token) QueryByPlatformMust(sqlx *sqlx.DB, platform int) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token_oauth2 where oauth2_id=? and expires_at > ? order by expires_at desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, platform, time.Now())
}

// UpdateToken update
func (a *ServerOAuth2Token) UpdateToken(sqlx *sqlx.DB) error {
	tic := sqlxc.TableIdxColumn{Table: TablePrefix + "token_oauth2", IDCol: "token_kid", IDVal: a.TokenID}
	if err := sqlxc.UpdateAndSaveByIDWithNamed(sqlx, nil, func() (string, map[string]interface{}, error) {
		return sqlxc.CreateUpdateSQLByNamedAndSkipNil(tic, a)
	}); err != nil {
		return err
	}
	return nil
}

// LockSync lock, 延迟锁定5秒
func (a *ServerOAuth2Token) LockSync(sqlx *sqlx.DB) error {
	SQL := "update {{TP}}oauth2_token set sync_lock=? where token_kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	_, err := sqlx.Exec(SQL, time.Now().Add(time.Duration(5)*time.Second), a.TokenID)
	return err
}
