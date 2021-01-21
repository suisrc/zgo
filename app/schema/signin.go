package schema

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/jmoiron/sqlx"
)

// SigninGpaUser user
type SigninGpaUser struct {
	ID     int        `db:"id" json:"-"`
	KID    string     `db:"kid" json:"id"`
	Name   string     `db:"name" json:"name"`
	Type   UserType   `db:"type" json:"-"`
	Status StatusType `db:"status" json:"-"`
}

// QueryByID sql 查询用户信息
func (a *SigninGpaUser) QueryByID(sqlx *sqlx.DB, id int, typ string) error {
	SQL := "select id, kid, name, status from {{TP}}user where id=?"
	params := []interface{}{id}
	if typ != "" {
		SQL += "and type=?"
		params = append(params, typ)
	}
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaOrgUser user
type SigninGpaOrgUser struct {
	ID       int            `tbl:"u" db:"id"`
	Type     UserType       `tbl:"u" db:"type"`
	UserID   int            `tbl:"t" db:"user_id"`
	OrgCode  string         `tbl:"t" db:"org_cod"`
	UnionKID string         `tbl:"t" db:"union_kid"`
	Name     string         `tbl:"t" db:"name"`
	CustomID sql.NullString `tbl:"t" db:"custom_id"`
	Status   StatusType     `tbl:"t" db:"status"`
}

// QueryByUserAndOrg sql select
func (a *SigninGpaOrgUser) QueryByUserAndOrg(sqlx *sqlx.DB, userid int, orgcode string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}tenant_user t inner join {{TP}}user u on u.id = t.user_id where t.user_id=? and t.org_cod=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, userid, orgcode)
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaAccount account
type SigninGpaAccount struct {
	ID           int            `db:"id"`
	PID          sql.NullInt64  `db:"pid"`           // 上级账户
	Account      string         `db:"account"`       // 账户
	AccountType  int            `db:"account_type"`  // 账户类型 1:name 2:mobile 3:email 4:openid 5:unionid 6:token
	PlatformKID  sql.NullString `db:"platform_kid"`  // 账户归属平台
	UserID       int            `db:"user_id"`       // 用户标识
	RoleID       sql.NullInt64  `db:"role_id"`       // 角色标识
	OrgCod       sql.NullString `db:"org_cod"`       // 角色标识
	Password     sql.NullString `db:"password"`      // 登录密码
	PasswordSalt sql.NullString `db:"password_salt"` // 密码盐值
	PasswordType sql.NullString `db:"password_type"` // 密码方式
	VerifySecret sql.NullString `db:"verify_secret"` // 校验密钥
	Status       StatusType     `db:"status"`        // 状态
	CreatedAt    sql.NullTime   `db:"created_at"`    // 创建时间
	UpdatedAt    sql.NullTime   `db:"updated_at"`    // 更新时间
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
}

// QueryByID 查询
func (a *SigninGpaAccount) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}account where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByAccount sql select
func (a *SigninGpaAccount) QueryByAccount(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("select " + sqlxc.SelectColumns(a))
	sqr.WriteString(" from {{TP}}account")
	sqr.WriteString(" where account=? and account_type=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		sqr.WriteString(" and platform_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and platform_kid is null")
	}
	sqr.WriteString(" and status=1")
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	// log.Println(SQL)
	return sqlx.Get(a, SQL, params...)
}

// QueryByParentAccount sql select
func (a *SigninGpaAccount) QueryByParentAccount(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	err := a.QueryByAccount(sqlx, acc, typ, kid)
	if err != nil {
		return err
	}
	if !a.PID.Valid {
		return errors.New("account pid is null")
	}
	paccount := SigninGpaAccount{}
	if err = paccount.QueryByID(sqlx, int(a.PID.Int64)); err != nil {
		return err
	} else if paccount.AccountType == int(AccountTypeName) || paccount.Status != StatusEnable {
		// 主账户不是密码账户或者主账户被禁用
		return errors.New("account pid is error")
	}
	// 使用主账户的密钥替换子账户
	a.Password = paccount.Password
	a.PasswordType = paccount.PasswordType
	a.PasswordSalt = paccount.PasswordSalt
	return nil
}

// QueryByAccountSkipStatus sql select
func (a *SigninGpaAccount) QueryByAccountSkipStatus(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("select " + sqlxc.SelectColumns(a))
	sqr.WriteString(" from {{TP}}account")
	sqr.WriteString(" where account=? and account_type=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		sqr.WriteString(" and platform_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and platform_kid is null")
	}
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// QueryByUserAndKind user and kid
func (a *SigninGpaAccount) QueryByUserAndKind(sqlx *sqlx.DB, uid int, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("select " + sqlxc.SelectColumns(a))
	sqr.WriteString(" from {{TP}}account")
	sqr.WriteString(" where user_id=? and account_type=?")

	params := []interface{}{uid, typ}
	if kid != "" {
		sqr.WriteString(" and platform_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and platform_kid is null")
	}
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// DeleteByUserAndKind user and kid
func (a *SigninGpaAccount) DeleteByUserAndKind(sqlx *sqlx.DB, uid int, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("delete from {{TP}}account")
	sqr.WriteString(" where user_id=? and account_type=?")

	params := []interface{}{uid, typ}
	if kid != "" {
		sqr.WriteString(" and platform_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and platform_kid is null")
	}
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)

	return sqlxc.DeleteOne(sqlx, SQL, params...)
}

// UpdateVerifySecret update verify secret
func (a *SigninGpaAccount) UpdateVerifySecret(sqlx *sqlx.DB) error {
	SQL := "update {{TP}}account set verify_secret=? where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	_, err := sqlx.Exec(SQL, a.VerifySecret.String, a.ID)
	return err
}

// UpdateAndSaveX update and save
func (a *SigninGpaAccount) UpdateAndSaveX(sqlx *sqlx.DB) error {
	Index := sqlxc.IdxColumn{Column: "id", ID: int64(a.ID)}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNil(TablePrefix+"account", Index, a)
	if err != nil {
		return err
	}

	res, err := sqlx.NamedExec(SQL, params)
	if err != nil {
		return err
	}
	if Index.ID == 0 {
		ID, _ := res.LastInsertId()
		a.ID = int(ID)
	}
	return nil
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaAccountToken account
type SigninGpaAccountToken struct {
	TokenID      string         `db:"token_kid"`
	AccountID    int            `db:"account_id"`
	DelayToken   sql.NullString `db:"delay_token"`
	DelayExpAt   sql.NullInt64  `db:"delay_exp"`
	OrgCode      sql.NullString `db:"org_cod"`
	AccessToken  sql.NullString `db:"access_token"`
	ExpiresAt    sql.NullInt64  `db:"expires_at"`
	RefreshToken sql.NullString `db:"refresh_token"`
	RefreshExpAt sql.NullInt64  `db:"refresh_exp"`
	CallCount    sql.NullInt64  `db:"call_count"`
	RefreshCount sql.NullInt64  `db:"refresh_count" set:"=refresh_count+1"`
	LastIP       sql.NullString `db:"last_ip"`
	LastAt       sql.NullTime   `db:"last_at"`
	ErrCode      sql.NullString `db:"error_code"`
	ErrMessage   sql.NullString `db:"error_message"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
	Number1      sql.NullInt64  `db:"number_1"` // 扩展字段
	String1      sql.NullString `db:"string_1"` // 扩展字段
}

// QueryByRefreshToken rtk
func (a *SigninGpaAccountToken) QueryByRefreshToken(sqlx *sqlx.DB, token string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where refresh_token=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, token)
}

// QueryByDelayToken rtk
func (a *SigninGpaAccountToken) QueryByDelayToken(sqlx *sqlx.DB, token string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where delay_token=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, token)
}

// QueryByTokenKID kid
func (a *SigninGpaAccountToken) QueryByTokenKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where token_kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByAccountAndClient kid
func (a *SigninGpaAccountToken) QueryByAccountAndClient(sqlx *sqlx.DB, accountID int, clientIP string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + " from {{TP}}token where account_id=?"
	params := []interface{}{accountID}
	if clientIP != "" {
		SQL += " and last_ip=?"
		params = append(params, clientIP)
	}
	SQL += " order by expires_at desc limit 1"

	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	log.Println(SQL)
	return sqlx.Get(a, SQL, params...)
}

// UpdateAndSaveByTokenKID 更新
func (a *SigninGpaAccountToken) UpdateAndSaveByTokenKID(sqlx *sqlx.DB, update bool) error {
	IDX := sqlxc.IdxColumn{Column: "token_kid", KID: a.TokenID, Create: !update, Update: update}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNilAndSet(TablePrefix+"token", IDX, a)
	if err != nil {
		return err
	}
	_, err = sqlx.NamedExec(SQL, params)
	return err
	// tx := sqlx.MustBegin()
	// tx.MustExec(SQL, params)
	// tx.Commit()
}

//=========================================================================
//=========================================================================
//=========================================================================

// type SigninGpaUserRole struct {
// 	UserID  int            `tbl:"ur" db:"user_id"`
// 	RoleID  int            `tbl:"ur" db:"role_id"`
// 	OrgCode string         `tbl:"ur" db:"org_cod"`
// 	OrgAdm  bool           `tbl:"ro" db:"org_adm"`
// 	KID     string         `tbl:"ro" db:"kid"`
// 	Name    string         `tbl:"ro" db:"name"`
// 	Status  StatusType     `tbl:"ro" db:"status"`
// 	SvcID   sql.NullInt64  `tbl:"ro" db:"svc_id"`
// 	SvcCode sql.NullString `tbl:"sv" db:"code"`
// }

// SigninGpaUserRole role
type SigninGpaUserRole struct {
	OrgAdm  bool           `tbl:"ro" db:"org_adm"`
	KID     string         `tbl:"ro" db:"kid"`
	Name    string         `tbl:"ro" db:"name"`
	SvcCode sql.NullString `tbl:"sv" db:"code"`
}

// QueryAllByUserID user -> user id / code -> org code
func (a *SigninGpaUserRole) QueryAllByUserID(sqlx *sqlx.DB, dest *[]SigninGpaUserRole, user int, code string) error {
	SQL := "select " + sqlxc.SelectColumns(a) + ` from {{TP}}user_role ur 
		inner join {{TP}}role ro on ro.id = ur.role_id 
		left  join {{TP}}app_service sv on sv.id = ro.svc_id 
		where ro.status = 1 and ur.user_id=?`
	params := []interface{}{user}
	if code != "" {
		SQL += " and ur.org_cod=?"
		params = append(params, code)
	} else {
		SQL += " and ur.org_cod is null"
	}
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL, params...)
}
