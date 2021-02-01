package module

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

// CasbinRule ...
type CasbinRule struct {
	Mid   int64  `db:"mid"`
	Ver   string `db:"ver"`
	PType string `db:"p_type"`
	V0    string `db:"v0"`
	V1    string `db:"v1"`
	V2    string `db:"v2"`
	V3    string `db:"v3"`
	V4    string `db:"v4"`
	V5    string `db:"v5"`
	V6    string `db:"v6"`
	V7    string `db:"v7"`
	V8    string `db:"v8"`
	V9    string `db:"v9"`
}

// Adapter 适配器
type Adapter struct {
	DB  *sqlx.DB // database
	Tbl string   //table name
	Mid int64    // model id
	Ver string   // model ver
}

// NewCasbinAdapter is the constructor for Adapter with existed connection
func NewCasbinAdapter(db *sqlx.DB, tbl string, mid int64, ver string) *Adapter {
	a := &Adapter{
		DB:  db,
		Tbl: tbl,
		Mid: mid,
		Ver: ver,
	}
	a.ensureTable()
	// runtime.SetFinalizer(a, finalizer)
	return a
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(m model.Model) error {
	var rules []CasbinRule
	err := a.DB.Select(&rules, fmt.Sprintf("SELECT * FROM `%s` WHERE mid = ? and ver = ?", a.Tbl), a.Mid, a.Ver)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		key := rule.PType
		sec := key[:1]
		args := []string{}

		v := reflect.ValueOf(rule)
		columnsIterator(func(c string, i int) (string, error) {
			if c[:1] == "v" {
				args = append(args, v.Field(i).String())
			}
			return "", nil
		}, rule)

		m[sec][key].Policy = append(m[sec][key].Policy, args)
		m[sec][key].PolicyMap[strings.Join(args, model.DefaultSep)] = len(m[sec][key].Policy) - 1
		// persist.LoadPolicyLine("", m)
	}
	return nil
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) (err error) {
	a.dropTable()
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := a.createPolicyRule(ptype, rule)
			err = a.insertPolicyLine(&line)
			if err != nil {
				return
			}
		}
	}
	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := a.createPolicyRule(ptype, rule)
			err = a.insertPolicyLine(&line)
			if err != nil {
				return
			}
		}
	}
	return
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) (err error) {
	line := a.createPolicyRule(ptype, rule)
	err = a.insertPolicyLine(&line)
	if err != nil {
		return
	}
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) (err error) {
	line := a.createPolicyRule(ptype, rule)
	err = a.deletePolicyLine(&line)
	if err != nil {
		return
	}
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (err error) {
	return
}

func (a *Adapter) rawDelete(line *CasbinRule) (err error) {
	args := []interface{}{}
	columns := strings.Builder{}

	v := reflect.ValueOf(line)
	columnsIterator(func(c string, i int) (string, error) {
		columns.WriteString("AND ")
		columns.WriteString(c)
		columns.WriteString(" = ?")
		args = append(args, v.Field(i).Interface())
		return "", nil
	}, line)

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s", a.Tbl, columns.String()[4:])
	_, err = a.DB.Exec(query, args...)
	if err != nil {
		return
	}
	return
}

func (a *Adapter) createPolicyRule(ptype string, rule []string) CasbinRule {
	line := CasbinRule{}
	line.Mid = a.Mid
	line.Ver = a.Ver
	line.PType = ptype

	v := reflect.ValueOf(line)
	columnsIterator(func(c string, i int) (string, error) {
		if c != "ver" && c[:1] == "v" {
			if idx, err := strconv.Atoi(c[1:]); err == nil {
				if idx < len(rule) {
					v.Field(i).SetString(rule[idx])
				} else {
					return "", errors.New("end assignment") // 结束遍历
				}
			}
		}
		return "", nil
	}, line)
	return line
}

func (a *Adapter) dropTable() {
	_, err := a.DB.Exec(fmt.Sprintf("DELETE FROM `%s`", a.Tbl))
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) ensureTable() {
	_, err := a.DB.Exec(fmt.Sprintf("SELECT 1 FROM `%s` LIMIT 1", a.Tbl))
	if err != nil {
		panic(err)
	}
}

func (a *Adapter) insertPolicyLine(line *CasbinRule) (err error) {
	columns := strings.Builder{}
	values := strings.Builder{}
	columnsItr(func(c string) string {
		columns.WriteString(",")
		columns.WriteString(c)
		values.WriteString(",:")
		values.WriteString(c)
		return ""
	}, line)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", a.Tbl, columns.String()[1:], values.String()[1:])
	_, err = a.DB.NamedExec(query, line)
	if err != nil {
		return
	}
	return
}

func (a *Adapter) deletePolicyLine(line *CasbinRule) (err error) {
	columns := strings.Builder{}
	columnsItr(func(c string) string {
		columns.WriteString("AND ")
		columns.WriteString(c)
		columns.WriteString(" = :")
		columns.WriteString(c)
		return ""
	}, line)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", a.Tbl, columns.String()[4:])
	_, err = a.DB.NamedExec(query, line)
	if err != nil {
		return
	}
	return
}

// columnsItr column
func columnsItr(iterater func(string) string, obj interface{}) ([]string, error) {
	return columnsIterator(func(c string, i int) (string, error) { return iterater(c), nil }, obj)
}

// columnsIterator column
func columnsIterator(iterater func(string, int) (string, error), obj interface{}) ([]string, error) {
	result := []string{}
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return result, nil
	}

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
		// if column == "" {
		// 	column = strings.ToLower(t.Field(i).Name)
		// }
		if iterater != nil {
			var err error
			column, err = iterater(column, i)
			if err != nil {
				return result, err
			}
		}
		if column != "" {
			result = append(result, column)
		}
	}
	return result, nil
}
