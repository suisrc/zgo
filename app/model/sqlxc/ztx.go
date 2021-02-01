package sqlxc

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TableIdxColumn id
// 1.Column="" 或者 Column = "id", 如果ID > 0 update, 否则 create
// 2.Update直接决定 update 或者 create
// 3.DB = nil 或者 Table = nil, 为 create
// 4.
type TableIdxColumn struct {
	Table  string
	IDCol  string      // index column
	IDVal  interface{} // index data
	Update sql.NullBool
	DB     *sqlx.DB
}

func (u *TableIdxColumn) isUpdate() bool {
	if u.IDCol == "" || u.IDCol == "id" {
		u.IDCol = "id"
		return u.IDVal.(int64) > 0
	} else if u.Update.Valid {
		return u.Update.Bool
	} else if u.DB == nil || u.Table == "" {
		return false
	}
	// 确定数据是否存在
	SQL := "select 1 from " + u.Table + " where " + u.IDCol + "=?"
	create := false
	if _, err := u.DB.Exec(SQL, u.IDVal); err != nil {
		create = IsNotFound(err)
	}
	u.Update = sql.NullBool{Valid: true, Bool: !create}
	return u.Update.Bool
}

// Index id
type Index struct {
	ID int64 `db:"id"`
}

// WithTx 执行带有事务的方法, 在一个事务中完成所有的内容
func WithTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx := db.MustBegin()
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}

// WithTxV 执行带有事务的方法, 在一个事务中完成所有的内容
func WithTxV(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) (interface{}, error)) (interface{}, error) {
	fnx := func(tx *sqlx.Tx, rr interface{}) (interface{}, error) {
		return fn(tx)
	}
	return WithTxVx(ctx, db, nil, fnx)
}

// WithTxVx 执行带有事务的方法, 在一个事务中完成所有的内容
func WithTxVx(ctx context.Context, db *sqlx.DB, rr interface{}, fn func(*sqlx.Tx, interface{}) (interface{}, error)) (interface{}, error) {
	tx := db.MustBegin()
	defer func() {
		if err := recover(); err != nil {
			// 发生中断异常
			tx.Rollback()
			panic(err)
		}
	}()
	res, err := fn(tx, rr)
	if err != nil {
		// 执行内容发生异常
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		// 提交发生异常
		return nil, errors.Wrapf(err, "committing transaction: %v", err)
	}
	return res, nil
}

// Struct2Map struct to map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		obj := v.Field(i).Interface()
		if idx := strings.Index(tag, ","); idx > 0 {
			tag = tag[:idx]
		}
		tag = strings.TrimSpace(tag)

		if tag == "" {
			data[strings.ToLower(t.Field(i).Name)] = obj
		} else {
			data[tag] = obj
		}
	}
	return data
}

// SelectColumns select column
func SelectColumns(obj interface{}) string {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	c := strings.Builder{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("db")
		if tag == "-" {
			continue
		}
		if idx := strings.Index(tag, ","); idx > 0 {
			tag = tag[:idx]
		}
		tag = strings.TrimSpace(tag)

		column := tag
		if column == "" {
			column = strings.ToLower(t.Field(i).Name)
		}
		c.WriteString(", ")

		if tbl := strings.TrimSpace(t.Field(i).Tag.Get("tbl")); tbl != "" && tbl != "-" {
			c.WriteString(tbl)
			if !strings.ContainsRune(tbl, '.') {
				c.WriteRune('.')
				c.WriteString(column)
			}
			c.WriteString(" as ")
		}
		c.WriteString(column)
	}
	return c.String()[1:]
}

// CreateUpdateSQLByNamedAndSkipNil create update sql by named
// 忽略空字段, 空字段不进行更新, 如果需要前置更新空字段, 直接使用 CreateUpdateSQLByNamed
func CreateUpdateSQLByNamedAndSkipNil(tic TableIdxColumn, obj interface{}) (string, map[string]interface{}, error) {
	return CreateUpdateSQLByNamed(tic, obj, func(t reflect.Type, n string, v interface{}, f *reflect.StructField) (interface{}, bool) {
		value := PickProxy(v)
		if n == "CreatedAt" && value == nil && !tic.isUpdate() || n == "UpdateAt" && value == nil {
			// 增加构建时间和更新时间字段
			value = NewNowTime(t)
		}
		if value == nil {
			return nil, false
		}
		return value, true
	}, nil)
}

// CreateUpdateSQLByNamedAndSkipNilAndSet create update sql by named
// 会处理 "set" 标签的内容 , 比如 `set:"=cloumn + 1"`
func CreateUpdateSQLByNamedAndSkipNilAndSet(tic TableIdxColumn, obj interface{}) (string, map[string]interface{}, error) {
	return CreateUpdateSQLByNamed(tic, obj,
		func(t reflect.Type, n string, v interface{}, f *reflect.StructField) (interface{}, bool) {
			value := PickProxy(v)
			if n == "CreatedAt" && value == nil && !tic.isUpdate() || n == "UpdatedAt" && value == nil {
				// 增加构建时间和更新时间字段
				value = NewNowTime(t)
			}
			if value == nil {
				if tic.isUpdate() && f.Tag.Get("set") != "" {
					return nil, true
				}
				return nil, false
			}
			return value, true
		}, func(c string, p map[string]interface{}, v interface{}, f *reflect.StructField) (string, bool) {
			tag := f.Tag.Get("set")
			if tag == "" {
				return "", false
			}
			if tag[:1] == "=" {
				return ", " + c + tag, true
			}
			return "", false
		})
}

// CreateUpdateSQLByNamed create update sql by named
func CreateUpdateSQLByNamed(tic TableIdxColumn, obj interface{},
	fix func(t reflect.Type, n string, v interface{}, f *reflect.StructField) (interface{}, bool),
	set func(c string, p map[string]interface{}, v interface{}, f *reflect.StructField) (string, bool),
) (string, map[string]interface{}, error) {

	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil, errors.New("obj is not struct")
	}

	SQL1 := strings.Builder{}
	SQL2 := strings.Builder{}
	params := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		key := t.Field(i).Name
		if key == "ID" {
			continue
		}
		field := t.Field(i)
		tag := field.Tag.Get("db")
		if tag == "-" {
			continue
		}
		if idx := strings.Index(tag, ","); idx > 0 {
			tag = tag[:idx]
		}
		tag = strings.TrimSpace(tag)
		column := tag
		if column == "" {
			column = strings.ToLower(key)
		}
		if tic.isUpdate() && column == tic.IDCol {
			continue
		}
		obj := v.Field(i).Interface()
		typ := t.Field(i).Type
		if fix != nil {
			res, ok := fix(typ, key, obj, &field)
			if !ok {
				continue
			}
			if res != obj {
				obj = res
			}
		} else {
			obj = PickProxy(obj)
		}
		if tic.isUpdate() {
			sql := ""
			ok := false
			if set != nil {
				sql, ok = set(column, params, obj, &field)
			}
			if ok {
				if sql != "" {
					SQL1.WriteString(sql)
				}
			} else {
				SQL1.WriteString(", " + column + "=:" + column)
				params[column] = obj
			}
		} else {
			SQL1.WriteString(", " + column)
			SQL2.WriteString(", :" + column)
			params[column] = obj
		}
	}

	if tic.isUpdate() {
		if SQL1.Len() == 0 {
			return "", nil, errors.New("no update values")
		}
		SQL := "update " + tic.Table + " set" + SQL1.String()[1:] + " where " + tic.IDCol + "=:" + tic.IDCol
		params[tic.IDCol] = tic.IDVal
		return SQL, params, nil
	}
	SQL := "insert into " + tic.Table + "(" + SQL1.String()[1:] + ") values (" + SQL2.String()[1:] + ")"
	return SQL, params, nil
}

// PickProxy 解除sql.NullXXX上的内容
func PickProxy(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	var b bool
	var v interface{}
	switch obj.(type) {
	case sql.NullBool:
		d := obj.(sql.NullBool)
		b, v = d.Valid, d.Bool
	case sql.NullFloat64:
		d := obj.(sql.NullFloat64)
		b, v = d.Valid, d.Float64
	case sql.NullInt32:
		d := obj.(sql.NullInt32)
		b, v = d.Valid, d.Int32
	case sql.NullInt64:
		d := obj.(sql.NullInt64)
		b, v = d.Valid, d.Int64
	case sql.NullString:
		d := obj.(sql.NullString)
		b, v = d.Valid, d.String
	case sql.NullTime:
		d := obj.(sql.NullTime)
		b, v = d.Valid, d.Time
	default:
		return obj
	}
	if b {
		return v
	}
	return nil
}

// NewNowTime 获取当前时间
func NewNowTime(t reflect.Type) interface{} {
	switch t.String() {
	case "sql.NullTime":
		return sql.NullTime{Valid: true, Time: time.Now()}
	case "time.Time":
		return time.Now()
	case "sql.NullInt64":
		return sql.NullInt64{Valid: true, Int64: time.Now().Unix()}
	case "int64":
		return time.Now().Unix()
	}
	return nil // 无法处理
}

// UpdateAndSaveByIDWithNamed update
func UpdateAndSaveByIDWithNamed(sqlx *sqlx.DB, id *Index, fn func() (string, map[string]interface{}, error)) error {
	SQL, params, err := fn()
	if err != nil {
		return err
	}
	res, err := sqlx.NamedExec(SQL, params)
	if err != nil {
		return err
	}
	if id != nil && id.ID == 0 {
		// 需要执行插入操作, 获取插入的ID
		if iid, err := res.LastInsertId(); err != nil {
			return err
		} else if iid > 0 {
			id.ID = iid
		}
	}
	return nil
}

// DeleteOne delete 1
func DeleteOne(sqlx *sqlx.DB, sql string, params ...interface{}) error {
	tx, err := sqlx.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(sql, params...)
	if err != nil {
		tx.Rollback()
		return err
	}
	rcc, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rcc != 1 {
		tx.Rollback()
		return errors.New("expected rows is one")
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
