package sqlxc

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/suisrc/zgo/app/model/gpa"
)

// NewClient client
func NewClient() (*sqlx.DB, func(), error) {
	//return NewClient2("sqlite3", "file:db1?mode=memory&cache=shared&_fk=1")
	return NewClient2(gpa.DatabaseType, gpa.DatabaseDSN())
}

// NewClient2 client
func NewClient2(typ, dsn string) (*sqlx.DB, func(), error) {
	db, err := sqlx.Connect(typ, dsn)
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxIdleConns(10)           // 最大空闲连接数
	db.SetMaxOpenConns(100)          // 数据库最大连接数
	db.SetConnMaxLifetime(time.Hour) //连接最长存活期，超过这个时间连接将不再被复用

	// run the auto migration tool.
	if gpa.TableSchemaInitSqlx || gpa.TableSchemaInit {
		if err := db.MustExec(schema); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		} else {
			// 防止其他持久化框架更新table结构
			gpa.TableSchemaInit = false
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
