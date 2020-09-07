package schema

import (
	"database/sql"
	"strings"

	"github.com/suisrc/zgo/app/model/sqlxc"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/modules/auth"
)

// SigninBody 登陆参数
type SigninBody struct {
	Username string `json:"username" binding:"required"` // 账户
	Password string `json:"password"`                    // 密码
	KID      string `json:"kid"`                         // 授权平台
	Client   string `json:"client"`                      // 子应用ID
	Captcha  string `json:"captcha"`                     // 验证码
	Code     string `json:"code"`                        // 标识码
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

// SigninQuery 登陆参数
type SigninQuery struct {
	Openid   string `form:"openid"`       // openid
	Code     string `form:"code"`         // code
	State    string `form:"state"`        // state
	Kid      string `form:"kid"`          // kid
	Redirect string `form:"redirect_uri"` // redirect_uri
}

// SigninResult 登陆返回值
type SigninResult struct {
	Status       string        `json:"status" default:"ok"`    // 'ok' | 'error' 不适用boolean类型是为了以后可以增加扩展
	Token        string        `json:"token,omitempty"`        // 令牌
	Expired      int64         `json:"expired,omitempty"`      // 过期时间
	RefreshToken string        `json:"refreshToken,omitempty"` // 刷新令牌
	Message      string        `json:"message,omitempty"`      // 消息,有限显示
	Roles        []interface{} `json:"roles,omitempty"`        // 多角色的时候，返回角色，重新确认登录
}

var _ auth.UserInfo = &SigninUser{}

// SigninUser 登陆用户信息
type SigninUser struct {
	UserName  string
	UserID    string
	RoleID    string
	TokenID   string
	Issuer    string
	Audience  string
	AccountID string
}

// GetUserName 用户名
func (s *SigninUser) GetUserName() string {
	return s.UserName
}

// GetUserID 用户ID
func (s *SigninUser) GetUserID() string {
	return s.UserID
}

// GetRoleID 角色ID
func (s *SigninUser) GetRoleID() string {
	return s.RoleID
}

// SetRoleID 角色ID
func (s *SigninUser) SetRoleID(nrole string) string {
	orole := s.RoleID
	s.RoleID = nrole
	return orole
}

// GetTokenID 令牌ID, 主要用于验证或者销毁令牌等关于令牌的操作
func (s *SigninUser) GetTokenID() string {
	return s.TokenID
}

// GetAccountID token
func (s *SigninUser) GetAccountID() string {
	return s.AccountID
}

// GetProps 获取私有属性,该内容会被加密, 注意:内容敏感,不要存储太多的内容
func (s *SigninUser) GetProps() (interface{}, bool) {
	return nil, false
}

// GetIssuer 令牌签发者
func (s *SigninUser) GetIssuer() string {
	return s.Issuer
}

// GetAudience 令牌接收者
func (s *SigninUser) GetAudience() string {
	return s.Audience
}

//==============================================================================

// SigninGpaUser user
type SigninGpaUser struct {
	ID     int    `db:"id" json:"-"`
	KID    string `db:"kid" json:"id"`
	Name   string `db:"name" json:"name"`
	Status bool   `db:"status" json:"-"`
}

// QueryByID sql select
func (a *SigninGpaUser) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select id, kid, name, status from {{TP}}user where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
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
	return sqlx.Select(dest, SQL, userid)
}

// SigninGpaAccount account
type SigninGpaAccount struct {
	ID           int            `db:"id"`
	PID          sql.NullInt64  `db:"pid"`
	Account      string         `db:"account"`
	AccountType  int            `db:"account_typ"`
	AccountKind  sql.NullString `db:"account_kid"`
	Password     sql.NullString `db:"password"`
	PasswordSalt sql.NullString `db:"password_salt"`
	PasswordType sql.NullString `db:"password_type"`
	VerifySecret sql.NullString `db:"verify_secret"`
	UserID       int            `db:"user_id"`
	RoleID       sql.NullInt64  `db:"role_id"`

	// SQLX1 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
	// SQLX2 int `sqlx:"from account where account=? and account_type='user' and platform='ZGO' and status=1"`
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
	sqr.WriteString(" where account=? and account_typ=?")

	params := []interface{}{acc, typ}
	if kid != "" {
		sqr.WriteString(" and account_kid=?")
		params = append(params, kid)
	} else {
		sqr.WriteString(" and account_kid is null")
	}
	sqr.WriteString(" and status=1")
	SQL := strings.ReplaceAll(sqr.String(), "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// UpdateVerifySecret update verify secret
func (a *SigninGpaAccount) UpdateVerifySecret(sqlx *sqlx.DB) error {
	SQL := "update {{TP}}account set verify_secret=? where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	_, err := sqlx.Exec(SQL, a.VerifySecret.String, a.ID)
	return err
}

// SigninGpaOAuth2Account account
type SigninGpaOAuth2Account struct {
	ID        int            `db:"id"`
	AccountID int            `db:"account_id"`
	ClientID  sql.NullInt64  `db:"client_id"`
	ClientKID sql.NullString `db:"client_kid"`
	UserKID   string         `db:"user_kid"`
	RoleKID   sql.NullString `db:"role_kid"`
	Expired   sql.NullInt64  `db:"expired"`
	LastIP    sql.NullString `db:"last_ip"`
	LastAt    sql.NullTime   `db:"last_at"`
	LimitExp  sql.NullTime   `db:"limit_exp"`
	LimitKey  sql.NullString `db:"limit_key"`
	Mode      sql.NullString `db:"mode"`
	Secret    sql.NullString `db:"secret"`
	Status    bool           `db:"status"`
}

// QueryByAccountAndClient kid
func (a *SigninGpaOAuth2Account) QueryByAccountAndClient(sqlx *sqlx.DB, accountID, clientID int) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_account where account_id=?"
	params := []interface{}{accountID}
	if clientID > 0 {
		SQL += " and client_id=?"
		params = append(params, clientID)
	} else {
		SQL += " and client_id is null"
	}
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// QueryByAccountAndClientK kid
func (a *SigninGpaOAuth2Account) QueryByAccountAndClientK(sqlx *sqlx.DB, accountID int, clientKID string) error {
	SQL := "select " + sqlxc.SelectColumns(a, "") + " from {{TP}}oauth2_account where account_id=?"
	params := []interface{}{accountID}
	if clientKID == "" {
		SQL += " and client_kid is null"
	} else {
		SQL += " and client_kid=?"
		params = append(params, clientKID)
	}
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, params...)
}

// UpdateAndSaveByAccountAndClient 更新
func (a *SigninGpaOAuth2Account) UpdateAndSaveByAccountAndClient(sqlx *sqlx.DB) (int64, error) {
	IDC := sqlxc.IDC{}
	if a.ClientKID.Valid {
		SQL := "select id from {{TP}}oauth2_account where account_id=? and client_kid=?"
		SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
		sqlx.Get(&IDC, SQL, a.AccountID, a.ClientKID)
	} else if a.ClientID.Valid {
		SQL := "select id from {{TP}}oauth2_account where account_id=? and client_id=?"
		SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
		sqlx.Get(&IDC, SQL, a.AccountID, a.ClientID)
	} else {
		SQL := "select id from {{TP}}oauth2_account where account_id=? and client_id is null"
		SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
		sqlx.Get(&IDC, SQL, a.AccountID)
	}
	SQL, params, err := sqlxc.CreateUpdateSQLByNamedAndSkipNil(TablePrefix+"oauth2_account", "id", IDC, a)
	if err != nil {
		return 0, err
	}
	// tx := sqlx.MustBegin()
	// tx.MustExec(SQL, params)
	// tx.Commit()
	res, err := sqlx.NamedExec(SQL, params)
	if err != nil {
		return 0, err
	}
	if IDC.ID > 0 {
		return IDC.ID, nil
	}
	return res.LastInsertId()
}
