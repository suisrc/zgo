package gpa

import (
	"github.com/jmoiron/sqlx"
)

// GPA golang persistence api 数据持久化
type GPA struct {
	Sqlx *sqlx.DB // sqlx client, 数据查询
}
