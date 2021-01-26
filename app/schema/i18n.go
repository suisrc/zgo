package schema

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// internationalization

// I18nGpaMessage message
type I18nGpaMessage struct {
	MID         string         `db:"mid"`
	Lang        string         `db:"lang"`
	Description sql.NullString `db:"description"`
	LeftDelim   sql.NullString `db:"left_delim"`
	RightDelim  sql.NullString `db:"right_delim"`
	Zero        sql.NullString `db:"zero"`
	One         sql.NullString `db:"one"`
	Few         sql.NullString `db:"few"`
	Many        sql.NullString `db:"many"`
	Other       sql.NullString `db:"other"`
}

// QueryAll 查询所有内容
func (a *I18nGpaMessage) QueryAll(sqlx *sqlx.DB, dest *[]I18nGpaMessage) error {
	SQL := "select id, mid, lang, description, left_delim, right_delim, zero, one, few, many, other from {{TP}}i18n_language where status=1"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL)
}

// I18nMessage 转换对象
func (a *I18nGpaMessage) I18nMessage() (*i18n.Message, language.Tag) {
	message := &i18n.Message{
		ID:          a.MID,
		Description: a.Description.String,
		LeftDelim:   a.LeftDelim.String,
		RightDelim:  a.RightDelim.String,
		Zero:        a.Zero.String,
		One:         a.One.String,
		Few:         a.Few.String,
		Many:        a.Many.String,
		Other:       a.Other.String,
	}

	tag := language.Make(a.Lang)
	return message, tag
}
