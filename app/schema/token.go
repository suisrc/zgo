package schema

import (
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

// TokenOAuth2 第三方登陆实体
type TokenOAuth2 struct {
	ID           int            `db:"id"`            // id
	PlatformID   int            `db:"oauth2_id"`     // 平台
	TokenID      sql.NullString `db:"token_kid"`     // 用户令牌标识
	AccessToken  sql.NullString `db:"access_token"`  // 访问令牌
	ExpiresIn    sql.NullInt64  `db:"expires_in"`    // 有限期间隔
	ExpiresTime  sql.NullTime   `db:"expires_time"`  // 凭据过期时间
	RefreshToken sql.NullString `db:"refresh_token"` // 刷新令牌
	RefreshCount sql.NullString `db:"refresh_count"` // 刷新次数
	SyncLock     sql.NullTime   `db:"sync_lock"`     // 同步锁
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	//TokenType    sql.NullString `db:"token_type"`
	//Version      sql.NullInt64  `db:"version" set:"=version+1"`
	//CallCount    sql.NullString `db:"call_count" set:"=version+1"`
}

// QueryByID kid
func (a *TokenOAuth2) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_token where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByPlatform platform
func (a *TokenOAuth2) QueryByPlatform(sqlx *sqlx.DB, platform int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_token where oauth2_id=? order by expires_time desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, platform)
}

// QueryByPlatformMust platform
func (a *TokenOAuth2) QueryByPlatformMust(sqlx *sqlx.DB, platform int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_token where oauth2_id=? and expires_time > ? order by expires_time desc limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, platform, time.Now())
}

// UpdateTokenOAuth2 update
func (a *TokenOAuth2) UpdateTokenOAuth2(sqlx *sqlx.DB) error {
	IDC := sqlxc.IDC{ID: int64(a.ID)}
	if err := sqlxc.UpdateAndSaveByIDWithNamed(sqlx, IDC, func() (string, map[string]interface{}, error) {
		return sqlxc.CreateUpdateSQLByNamedAndSkipNil(TablePrefix+"oauth2_token", "id", IDC, a)
	}); err != nil {
		return err
	} else if a.ID == 0 {
		a.ID = int(IDC.ID)
	}
	return nil
}

// LockSync lock, 延迟锁定5秒
func (a *TokenOAuth2) LockSync(sqlx *sqlx.DB) error {
	SQL := "update {{TP}}oauth2_token set sync_lock=? where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	_, err := sqlx.Exec(SQL, time.Now().Add(time.Duration(5)*time.Second), a.ID)
	return err
}
