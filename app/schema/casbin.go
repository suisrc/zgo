package schema

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

// CasbinGpaSvcAud 企业授权服务
type CasbinGpaSvcAud struct {
	ID       int            `tbl:"sa" db:"id"`
	SvcID    sql.NullInt64  `tbl:"sa" db:"svc_id"`
	SvcCode  sql.NullString `tbl:"sv" db:"code"`
	OrgCode  sql.NullString `tbl:"sa" db:"org_cod"`
	Audience sql.NullString `tbl:"sa" db:"audience"`
	Resource sql.NullString `tbl:"sa" db:"resource"`
}

// CasbinGpaSvcAudSlice slice
type CasbinGpaSvcAudSlice []CasbinGpaSvcAud

func (p CasbinGpaSvcAudSlice) Len() int           { return len(p) }
func (p CasbinGpaSvcAudSlice) Less(i, j int) bool { return p[i].Resource.String > p[j].Resource.String }
func (p CasbinGpaSvcAudSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// QueryByAudAndResAndOrg 查询, 1st: aud + res, 2rd: res, 3th: aud, 如果三种方式都无法找到， 确定发生异常
func (a *CasbinGpaSvcAud) QueryByAudAndResAndOrg(sqlx *sqlx.DB, aud, res, org string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}app_service_audience sa inner join app_service sv on sv.id = sa.svc_id where "
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	SQLX := []string{
		"sa.audience = ? and sa.resource = ? limit 1",
		"sa.audience = ? and sa.resource is null or sa.audience is null and sa.resource = ? desc sa.resource limit 1",
	}
	for _, sx := range SQLX {
		if err := sqlx.Get(a, SQL+sx, aud, res); err != nil {
			if !sqlxc.IsNotFound(err) {
				return err
			}
		} else {
			break
		}
	}
	return errors.New("sql: no rows in result set")
}

// CasbinGpaSvcOrg 授权给机构的服务
type CasbinGpaSvcOrg struct {
	OrgCode string       `db:"org_cod"`
	SvcID   int          `db:"svc_id"`
	Expired sql.NullTime `db:"expired"`
	Status  StatusType   `db:"status"`
}

// QueryByOrgAndSvc 查询
func (a *CasbinGpaSvcOrg) QueryByOrgAndSvc(sqlx *sqlx.DB, org string, sid int) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}app_service_org where org_cod = ? and svc_id = ?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, org, sid)
}

//=============================================================================================
//=============================================================================================
//=============================================================================================

// CasbinGpaModel model
type CasbinGpaModel struct {
	ID        int            `db:"id"`
	OrgCode   sql.NullString `db:"org_cod"`
	Name      string         `db:"name"`
	Statement string         `db:"statement"`
}

// QueryAllByOrg 查询
func (a *CasbinGpaModel) QueryAllByOrg(sqlx *sqlx.DB, org string) (*[]CasbinGpaModel, error) {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}policy_model pm where pm.org_cod is null or pm.org_cod = ?`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	res := []CasbinGpaModel{}
	if err := sqlx.Select(&res, SQL, org); err != nil {
		return nil, err
	}
	return &res, nil
}

// QueryByOrg 查询
func (a *CasbinGpaModel) QueryByOrg(sqlx *sqlx.DB, org string) (*CasbinGpaModel, error) {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}policy_model pm 
			where pm.org_cod is null or pm.org_cod = ? 
			desc pm.org_cod, pm.id limit 1`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	if err := sqlx.Get(a, SQL, org); err != nil {
		return nil, err
	}
	return a, nil
}

//==========================================================================================
//==========================================================================================
//==========================================================================================

// CasbinGpaRoleRole RoleRole
type CasbinGpaRoleRole struct {
	ParentName string         `tbl:"rp.name" db:"pn_name"`
	ParentSvc  sql.NullString `tbl:"sp.name" db:"ps_name"`
	ChildName  string         `tbl:"rc.name" db:"cn_name"`
	ChildSvc   sql.NullString `tbl:"sc.name" db:"cs_name"`
	OrgCode    sql.NullString `tbl:"rr" db:"org_cod"`
}

// QueryByOrg 查询 有效状态为1
func (a *CasbinGpaRoleRole) QueryByOrg(sqlx *sqlx.DB, org string) (*[]CasbinGpaRoleRole, error) {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}role_role rr 
		inner join {{TP}}role rp on rp.id = rr.pid
		inner join {{TP}}role rc on rc.id = rr.cid
		left join {{TP}}app_service sp on sp.id = rp.svc_id
		left join {{TP}}app_service sc on sc.id = rc.svc_id
		where (rr.org_cod is null or rr.org_cod = ?) and rp.status = 1 and rc.status = 1`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	res := []CasbinGpaRoleRole{}
	err := sqlx.Select(&res, SQL, org)
	return &res, err
}

// CasbinGpaRolePolicy RolePolicy
type CasbinGpaRolePolicy struct {
	Role    string         `tbl:"ro.name" db:"r_name"`
	Svc     sql.NullString `tbl:"rv.name" db:"v_name"`
	Policy  string         `tbl:"po.name" db:"p_name"`
	OrgCode sql.NullString `tbl:"rp" db:"org_cod"`
}

// QueryByOrg 查询 有效状态为1
func (a *CasbinGpaRolePolicy) QueryByOrg(sqlx *sqlx.DB, org string) (*[]CasbinGpaRolePolicy, error) {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}role_policy rp 
		inner join {{TP}}role ro on ro.id = rp.role_id
		inner join {{TP}}policy po on po.id = rp.plcy_id
		left join {{TP}}app_service sv on sv.id = ro.svc_id
		where (rp.org_cod is null or rp.org_cod = ?) and ro.status = 1 and po.status = 1`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	res := []CasbinGpaRolePolicy{}
	if err := sqlx.Select(&res, SQL, org); res != nil {
		return nil, err
	}
	return &res, nil
}

// CasbinGpaPolicyStatement PolicyStatement 执行策略
type CasbinGpaPolicyStatement struct {
	Name      string         `tbl:"po.name" db:"p_name"`
	OrgCode   sql.NullString `tbl:"po" db:"org_cod"`
	Status    StatusType     `tbl:"po" db:"status"`
	Version   int            `tbl:"po" db:"version"` // 版本必须匹配
	Effect    bool           `tbl:"ps" db:"effect"`
	Action    sql.NullString `tbl:"ps" db:"action"`
	Resource  sql.NullString `tbl:"ps" db:"resource"`
	Condition sql.NullString `tbl:"ps" db:"condition"`
}

// QueryByOrg 查询 有效状态为1
func (a *CasbinGpaPolicyStatement) QueryByOrg(sqlx *sqlx.DB, org string) (*[]CasbinGpaPolicyStatement, error) {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}policy_statement ps 
		inner join {{TP}}policy py on py.id = ps.plcy_id and py.version = ps.version
		where (py.org_cod is null or py.org_cod = ?) and py.status = 1`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	res := []CasbinGpaPolicyStatement{}
	if err := sqlx.Select(&res, SQL, org); err != nil {
		return nil, err
	}
	return &res, nil
}

// CasbinGpaPolicyServiceAction PolicyServiceAction
type CasbinGpaPolicyServiceAction struct {
	Name     string         `tbl:"psa.name" db:"a_name"`
	SvcCode  string         `tbl:"psv.code" db:"v_code"`
	Resource sql.NullString `tbl:"psa" db:"resource"`
	Status   StatusType     `tbl:"psa" db:"status"`
}

// QueryActionByNameAndSvc 查询相应时间 * -> % / ? -> _
func (a *CasbinGpaPolicyServiceAction) QueryActionByNameAndSvc(sqlx *sqlx.DB, name, svc string) (*[]CasbinGpaPolicyServiceAction, error) {
	params := []interface{}{svc, strings.ReplaceAll(strings.ReplaceAll(name, "?", "_"), "*", "%")}
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}policy_service_action psa 
		inner join app_service psv on psv.id = psa.svc_id and psv.code = ? 
		where pas.name like ?`
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	res := []CasbinGpaPolicyServiceAction{}
	if err := sqlx.Select(&res, SQL, params...); err != nil {
		return nil, err
	}
	return &res, nil
}
