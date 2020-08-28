package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// JwtGpaOpts jwt配置信息
type JwtGpaOpts struct {
	ID          string         `db:"id"`
	KID         string         `db:"kid"`
	Expired     int            `db:"expired"`
	Secret      string         `db:"token_secret"`
	Audience    sql.NullString `db:"audience"`
	Issuer      sql.NullString `db:"issuer"`
	SigninURL   sql.NullString `db:"signin_url"`
	SigninForce sql.NullBool   `db:"signin_force"`
	SigninCheck sql.NullBool   `db:"signin_check"`
	SecretByte  []byte         `db:"-"`
}

// QueryAll sql select
func (*JwtGpaOpts) QueryAll(sqlx *sqlx.DB, dest *[]JwtGpaOpts) error {
	SQL := "select id, kid, expired, audience, issuer, token_secret, signin_url, signin_force, signin_check from {{TP}}oauth2_client where status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// QueryByKID kid
func (a *JwtGpaOpts) QueryByKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select id, kid, expired, audience, issuer, token_secret, signin_url, signin_force, signin_check from {{TP}}oauth2_client where kid=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByAudience kid
func (a *JwtGpaOpts) QueryByAudience(sqlx *sqlx.DB, audience string) error {
	SQL := "select id, kid, expired, audience, issuer, token_secret, signin_url, signin_force, signin_check from {{TP}}oauth2_client where audience=? and status=1 limit 1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, audience)
}
