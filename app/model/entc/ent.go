package entc

import (
	"context"
	"log"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/suisrc/zgo/app/model/ent"
	"github.com/suisrc/zgo/app/model/gpa"
)

// NewClient client
func NewClient() (*ent.Client, func(), error) {

	// drv, err := sql.Open("sqlite3", "file:db1?mode=memory&cache=shared&_fk=1")
	drv, err := sql.Open(gpa.DatabaseType, gpa.DatabaseDSN())
	if err != nil {
		return nil, nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	client := ent.NewClient(ent.Driver(drv))
	//client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	//if err != nil {
	//	return nil, nil, err
	//}

	// run the auto migration tool.
	if gpa.TableSchemaInitEnt || gpa.TableSchemaInit {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		} else {
			// 防止其他持久化框架更新table结构
			gpa.TableSchemaInit = false
		}
	}

	// defer client.Close()
	// runtime.SetFinalizer(client, func(client *Client){client.Close()})
	clean := func() {
		client.Close()
	}
	return client, clean, nil
}
