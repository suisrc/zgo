package sqlxc

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	zdbc "github.com/suisrc/zgo/app/model"
)

// NewClient client
func NewClient() (*sqlx.DB, func(), error) {
	//db, err := sqlx.Connect("sqlite3", "file:db1?mode=memory&cache=shared&_fk=1")
	db, err := sqlx.Connect(zdbc.DatabaseType, zdbc.DatabaseDSN())
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxIdleConns(10)           // 最大空闲连接数
	db.SetMaxOpenConns(100)          // 数据库最大连接数
	db.SetConnMaxLifetime(time.Hour) //连接最长存活期，超过这个时间连接将不再被复用

	// run the auto migration tool.
	if zdbc.TableSchemaInitSqlx || zdbc.TableSchemaInit {
		if err := db.MustExec(schema); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		} else {
			// 防止其他持久化框架更新table结构
			zdbc.TableSchemaInit = false
		}
	}

	// defer client.Close()
	clean := func() {
		db.Close()
	}
	return db, clean, nil
}

// IsNotFound of sqlx
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return "sql: no rows in result set" == err.Error()
}
