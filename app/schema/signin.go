package schema

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/jmoiron/sqlx"
)

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password"`                    // 密码
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
	KID      string `json:"kid"`                         // 授权平台
	Org      string `json:"org"`                         // 租户
	Role     string `json:"role"`                        // 角色
	Domain   string `json:"host"`                        // 域, 如果无,使用c.Reqest.Host代替
}

// SigninOfCaptcha 使用登陆发生认证信息
type SigninOfCaptcha struct {
	Mobile string `form:"mobile"` // 手机
	Email  string `form:"email"`  // 邮箱
	Openid string `form:"openid"` // openid
	KID    string `form:"kid"`    // 平台标识
}

// SigninOfOAuth2 登陆参数
type SigninOfOAuth2 struct {
	Code     string `form:"code"`     // 票据
	State    string `form:"state"`    // 验签
	Scope    string `form:"scope"`    // 作用域
	KID      string `form:"kid"`      // kid
	Org      string `form:"org"`      // 租户
	Role     string `form:"role"`     // 角色
	Domain   string `form:"host"`     // 域, 如果无,使用c.Reqest.Host代替
	Redirect string `form:"redirect"` // redirect
}

// SigninResult 登陆返回值
type SigninResult struct {
	TokenStatus  string        `json:"status" default:"ok"`                   // 'ok' | 'error' 不适用boolean类型是为了以后可以增加扩展
	TokenID      string        `json:"token_id,omitempty"`                    // 访问令牌ID
	AccessToken  string        `json:"access_token,omitempty"`                // 访问令牌
	TokenType    string        `json:"token_type,omitempty" default:"bearer"` // 令牌类型
	ExpiresAt    int64         `json:"expires_at,omitempty"`                  // 过期时间
	ExpiresIn    int64         `json:"expires_in,omitempty"`                  // 过期时间
	RefreshToken string        `json:"refresh_token,omitempty"`               // 刷新令牌
	RefreshExpAt int64         `json:"refresh_expires,omitempty"`             // 刷新令牌过期时间
	Redirect     string        `json:"redirect_uri,omitempty"`                // redirect_uri
	Message      string        `json:"message,omitempty"`                     // 消息,有限显示 // Message 和 Datas 一般用户发生异常后回显
	Params       []interface{} `json:"params,omitempty"`                      // 多租户多角色的时候，返回角色，重新确认登录
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaUser user
type SigninGpaUser struct {
	ID     int        `db:"id" json:"-"`
	KID    string     `db:"kid" json:"id"`
	Name   string     `db:"name" json:"name"`
	Type   string     `db:"type" json:"-"`
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

// SigninGpaRole role
type SigninGpaRole struct {
	ID     int     `db:"id" json:"-"`
	KID    string  `db:"kid" json:"id"`
	Name   string  `db:"name" json:"name"`
	Domain *string `db:"domain" json:"domain"`
	//Domain sql.NullString `db:"domain" json:"domain"`
}

// QueryByID sql select
func (a *SigninGpaRole) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name, domain from {{TP}}role where id=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByKID sql select
func (a *SigninGpaRole) QueryByKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select id, kid, name, domain from {{TP}}role where kid=? and status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByUserID sql select
func (a *SigninGpaRole) QueryByUserID(sqlx *sqlx.DB, dest *[]SigninGpaRole, userid int) error {
	SQL := "select r.id, r.kid, r.name, r.domain from {{TP}}user_role ur inner join {{TP}}role r on r.id=ur.role_id where ur.user_id=? and r.status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	log.Println(SQL)
	return sqlx.Select(dest, SQL, userid)
}

//=========================================================================
//=========================================================================
//=========================================================================

// SigninGpaOrgUser user
type SigninGpaOrgUser struct {
	UserID   int            `db:"user_id"`
	OrgCode  string         `db:"org_cod"`
	UnionKID string         `db:"union_kid"`
	Name     string         `db:"name"`
	CustomID sql.NullString `db:"custom_id"`
	Status   StatusType     `db:"status"`
}

// QueryByUserAndOrg sql select
func (a *SigninGpaOrgUser) QueryByUserAndOrg(sqlx *sqlx.DB, userid int, orgcode string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}tenant_user where user_id=? and org_cod=?"
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
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}account where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByAccount sql select
func (a *SigninGpaAccount) QueryByAccount(sqlx *sqlx.DB, acc string, typ int, kid string) error {
	sqr := strings.Builder{}
	sqr.WriteString("select " + sqlxc.SelectColumns(a, ""))
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
	sqr.WriteString("select " + sqlxc.SelectColumns(a, ""))
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
	sqr.WriteString("select " + sqlxc.SelectColumns(a, ""))
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
	AccessToken  sql.NullString `db:"access_token"`
	ExpiresAt    sql.NullInt64  `db:"expires_at"`
	RefreshToken sql.NullString `db:"refresh_token"`
	RefreshExpAt sql.NullInt64  `db:"refresh_expires"`
	CallCount    sql.NullInt64  `db:"call_count"`
	RefreshCount sql.NullInt64  `db:"refresh_count" set:"=refresh_count+1"`
	LastIP       sql.NullString `db:"last_ip"`
	LastAt       sql.NullTime   `db:"last_at"`
	ErrCode      sql.NullString `db:"error_code"`
	ErrMessage   sql.NullString `db:"error_message"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=version+1"`
}

// QueryByRefreshToken rtk
func (a *SigninGpaAccountToken) QueryByRefreshToken(sqlx *sqlx.DB, token string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}token where refresh_token=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, token)
}

// QueryByTokenKID kid
func (a *SigninGpaAccountToken) QueryByTokenKID(sqlx *sqlx.DB, kid string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}token where token_kid=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, kid)
}

// QueryByAccountAndClient kid
func (a *SigninGpaAccountToken) QueryByAccountAndClient(sqlx *sqlx.DB, accountID int, clientIP string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}token where account_id=?"
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
