package zdbc

import (
	"github.com/suisrc/zgo/modules/config"
	// 引入数据库
	// _ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

var (
	// sqlite3

	// DatabaseType type
	// DatabaseType = "sqlite3"
	// DatabaseDSN dsn
	// DatabaseDSN = "file:db1?mode=memory&cache=shared&_fk=1"

	// mysql

	// DatabaseType type
	DatabaseType = "mysql"
	// DatabaseDSN dsn
	DatabaseDSN = func() string {
		return config.C.MySQL.DSN()
	}
)

var (
	// TableSchemaInit 是否初始化数据表结构
	TableSchemaInit = false

	// TableSchemaInitEnt 强制使用ent更新表结构, TableSchemaInit => true 无效
	TableSchemaInitEnt = false

	// TableSchemaInitSqlx 强制使用sqlx更新表结构, TableSchemaInit => true 无效
	TableSchemaInitSqlx = false
)
