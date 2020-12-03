package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CasbinGpaResource policy => p
type CasbinGpaResource struct {
	ID       int            `db:"id"`
	Resource sql.NullString `db:"resource"`
	Domain   sql.NullString `db:"domain"`
	Methods  sql.NullString `db:"methods"`
	Path     sql.NullString `db:"path"`
	Netmask  sql.NullString `db:"netmask"`
	Allow    sql.NullBool   `db:"allow"`
	// description sql.NullString `db:"description"`
	// status      sql.NullBool   `db:"status"`
}

// QueryAll sql select
func (*CasbinGpaResource) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaResource) error {
	SQL := "select id, resource, methods, path, netmask, allow from {{TP}}resource where status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// CasbinGpaResourceRole policy => g
type CasbinGpaResourceRole struct {
	ID       int            `db:"id"`
	Role     sql.NullString `db:"role"`
	Resource sql.NullString `db:"resource"`
}

// QueryAll sql select
func (*CasbinGpaResourceRole) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaResourceRole) error {
	SQL := "select rr.id, r.kid, rr.resource as role from {{TP}}resource_role rr inner join {{TP}}role r on r.id = rr.role_id where r.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// CasbinGpaResourceUser policy => g
type CasbinGpaResourceUser struct {
	ID       int            `db:"id"`
	User     sql.NullString `db:"user"`
	Resource sql.NullString `db:"resource"`
}

// QueryAll sql select
func (*CasbinGpaResourceUser) QueryAll(sqlx *sqlx.DB, dest *[]CasbinGpaResourceUser) error {
	SQL := "select ru.id, u.kid as user, ru.resource from {{TP}}resource_user ru inner join {{TP}}user u on u.id = ru.user_id where u.status=1"
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
