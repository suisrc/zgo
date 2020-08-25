package schema

import "database/sql"

// JwtOpts jwt配置信息
type JwtOpts struct {
	ID          string         `db:"id"`
	KID         string         `db:"kid"`
	Expired     int            `db:"expired"`
	Secret      string         `db:"token_secret"`
	Audience    sql.NullString `db:"audience"`
	Issuer      sql.NullString `db:"issuer"`
	SigninURL   sql.NullString `db:"signin_url"`
	SigninForce sql.NullBool   `db:"signin_force"`
	SigninCheck sql.NullBool   `db:"signin_check"`
}

// SQLByAll sql select
func (*JwtOpts) SQLByAll() string {
	return "select id, kid, expired, audience, issuer, token_secret, signin_url, signin_force, signin_check from user where status=1"
}
