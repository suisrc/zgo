package gpa

import (
	"github.com/jmoiron/sqlx"
)

// GPA golang persistence api 数据持久化
type GPA struct {
	// Entc *ent.Client // ent client, 数据修改和插入
	Sqlx *sqlx.DB // sqlx client, 数据查询
}
