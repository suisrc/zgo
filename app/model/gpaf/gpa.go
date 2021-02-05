package gpaf

import (
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/modules/config"
)

// NewGPA client
func NewGPA() (gpa.GPA, func(), error) {
	db1, clean1, err1 := sqlxc.NewClient()
	if err1 != nil {
		return gpa.GPA{}, nil, err1
	}
	if config.C.MySQL2.Host == "" {
		// 第二数据源没有配置, 使用第一数据源代替
		clean := func() {
			clean1()
		}
		gpa := gpa.GPA{
			Sqlx:  db1,
			Sqlx2: db1,
		}
		return gpa, clean, nil
	}
	db2, clean2, err2 := sqlxc.NewClient2(gpa.DatabaseType, gpa.Database2DSN())
	if err2 != nil {
		return gpa.GPA{}, nil, err1
	}

	clean := func() {
		clean1()
		clean2()
	}
	gpa := gpa.GPA{
		Sqlx:  db1,
		Sqlx2: db2,
	}
	return gpa, clean, nil
}
