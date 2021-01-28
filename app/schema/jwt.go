package schema

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// JwtGpaOpts jwt配置信息
type JwtGpaOpts struct {
	ID          int            `db:"id"`
	KID         string         `db:"kid"`
	Secret      string         `db:"secret"`
	Expired     int            `db:"expired"`
	Refresh     int            `db:"refresh"`
	Issuer      sql.NullString `db:"issuer"`
	Audience    sql.NullString `db:"audience"`
	SigninURL   sql.NullString `db:"signin_url"`
	SigninForce sql.NullBool   `db:"signin_force"`
	SigninCheck sql.NullBool   `db:"signin_check"`
	SecretByte  []byte         `db:"-"`
}

// QueryAll sql select
func (*JwtGpaOpts) QueryAll(sqlx *sqlx.DB) (*[]JwtGpaOpts, error) {
	return nil, nil
}

// QueryByKID kid
func (a *JwtGpaOpts) QueryByKID(sqlx *sqlx.DB, kid string) error {
	return nil
}

// QueryByAudience kid
func (a *JwtGpaOpts) QueryByAudience(sqlx *sqlx.DB, audience string, org string) error {
	return nil
}
