package gpa

import (
	"github.com/jmoiron/sqlx"
)

// entc.NewClient,                 // 数据库连接注册
// sqlxc.NewClient,                // 数据库连接注册
// wire.Struct(new(gpa.GPA), "*"), // 数据库服务

// GPA golang persistence api 数据持久化
type GPA struct {
	// Entc *ent.Client // ent client, 数据修改和插入
	Sqlx  *sqlx.DB // sqlx client, 数据查询
	Sqlx2 *sqlx.DB // sqlx client, 数据查询
}
