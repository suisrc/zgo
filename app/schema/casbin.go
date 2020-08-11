package schema

import "database/sql"

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

// SQLByALL sql select
func (*CasbinGpaResource) SQLByALL() string {
	return "select id, resource, methods, path, netmask, allow from resource where status=1"
}

// CasbinGpaResourceRole policy => g
type CasbinGpaResourceRole struct {
	ID       int            `db:"id"`
	Role     sql.NullString `db:"role"`
	Resource sql.NullString `db:"resource"`
}

// SQLByALL sql select
func (*CasbinGpaResourceRole) SQLByALL() string {
	return "select rr.id, rr.resource, r.uid as role from resource_role rr inner join role r on r.id = rr.role_id where r.status=1"
}

// CasbinGpaRoleRole policy => g
type CasbinGpaRoleRole struct {
	ID    int            `db:"id"`
	Owner sql.NullString `db:"owner"`
	Child sql.NullString `db:"child"`
}

// SQLByALL sql select
func (*CasbinGpaRoleRole) SQLByALL() string {
	return `select rr.id, ro.uid as owner, rc.uid as child 
			from role_role rr 
			inner join role ro on ro.id = rr.owner_id 
			inner join role rc on rc.id = rr.child_id 
			where ro.status=1 and rc.status=1`
}
