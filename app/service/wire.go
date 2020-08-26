package service

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/ent"
	"github.com/suisrc/zgo/app/model/entc"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/modules/passwd"
)

// ServiceSet wire注入服务
var ServiceSet = wire.NewSet(
	NewAuther,                               // Auther注册
	wire.Struct(new(passwd.Validator), "*"), // 密码验证
	entc.NewClient,                          // 数据库连接注册
	sqlxc.NewClient,                         // 数据库连接注册
	wire.Struct(new(GPA), "*"),              // 数据库服务
	wire.Struct(new(AuthOpts), "GPA"),       // Auther依赖
	// 服务
	wire.Struct(new(Signin), "*"),
)

//======================================
// 分割线
//======================================

// ResultRef 返回值暂存器
type ResultRef struct {
	D interface{}
}

// GPA golang persistence api 数据持久化
type GPA struct {
	Entc *ent.Client // ent client, 数据修改和插入
	Sqlx *sqlx.DB    // sqlx client, 数据查询
}
