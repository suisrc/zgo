package schema

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// DemoSet 示例对象
type DemoSet struct {
	Code string `json:"code" binding:"required"` // 编号
	Name string `json:"name" binding:"required"` // 名称
	Memo string `json:"memo"`                    // 备注
}

// QueryByID select by id
func (a *DemoSet) QueryByID(sqlx *sqlx.DB, id int) error {
	SQL := "select code, name, memo from {{TP}}demo where id=?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Get(a, SQL, id)
}

// QueryByNamelike select by name like
func (a *DemoSet) QueryByNamelike(sqlx *sqlx.DB, dest *[]DemoSet, name string) error {
	SQL := "select code, name, memo from {{TP}}demo where name like ?"
	SQL = strings.ReplaceAll(SQL, "{{TP}}", TablePrefix)
	return sqlx.Select(dest, SQL, "%"+name+"%")
}
