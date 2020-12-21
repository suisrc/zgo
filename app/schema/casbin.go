package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CasbinGpaGateway policy => p
type CasbinGpaGateway struct {
	ID      int            `db:"id"`
	Name    sql.NullString `db:"name"`
	Domain  sql.NullString `db:"domain"`
	Methods sql.NullString `db:"methods"`
	Path    sql.NullString `db:"path"`
	Netmask sql.NullString `db:"netmask"`
	Allow   sql.NullBool   `db:"allow"`
	// description sql.NullString `db:"description"`
	// status      sql.NullBool   `db:"status"`
}

// QueryAll sql select
func (*CasbinGpaGateway) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaGateway) error {
	SQL := "select id, name, methods, path, netmask, allow from {{TP}}gateway where status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// CasbinGpaGatewayRole policy => g
type CasbinGpaRoleGateway struct {
	ID      int            `db:"id"`
	Role    sql.NullString `db:"role"`
	Gateway sql.NullString `db:"gateway"`
}

// QueryAll sql select
func (*CasbinGpaRoleGateway) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaRoleGateway) error {
	SQL := "select rg.id, r.kid as role, rg.gateway from {{TP}}role_gateway rg inner join {{TP}}role r on r.id = rg.role_id where r.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// CasbinGpaGatewayUser policy => g
type CasbinGpaUserGateway struct {
	ID      int            `db:"id"`
	User    sql.NullString `db:"user"`
	Gateway sql.NullString `db:"gateway"`
}

// QueryAll sql select
func (*CasbinGpaUserGateway) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaUserGateway) error {
	SQL := "select ug.id, u.kid as user, ug.gateway from {{TP}}user_gateway ug inner join {{TP}}user u on u.id = ug.user_id where u.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// CasbinGpaRoleRole policy => g
type CasbinGpaRoleRole struct {
	ID    int            `db:"id"`
	Owner sql.NullString `db:"owner"`
	Child sql.NullString `db:"child"`
}

// QueryAll sql select
func (*CasbinGpaRoleRole) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaRoleRole) error {
	SQL := `select rr.id, ro.kid as owner, rc.kid as child 
			from {{TP}}role_role rr 
			inner join {{TP}}role ro on ro.id = rr.owner_id 
			inner join {{TP}}role rc on rc.id = rr.child_id 
			where ro.status=1 and rc.status=1`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}
