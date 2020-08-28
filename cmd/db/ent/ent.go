package entcmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/suisrc/zgo/modules/logger"
)

//ModelFile file
type ModelFile struct {
	Model  string
	Output string
}

// RunBuild build
func (a *ModelFile) RunBuild() error {
	log.Println("starting build mysql script")
	log.Println(fmt.Sprintf("model input: %s", a.Model))
	log.Println(fmt.Sprintf("model output: %s", a.Output))

	bytes, err := ioutil.ReadFile(a.Model)
	if err != nil {
		log.Println(fmt.Sprintf("model file error: %s", logger.ErrorWW(err)))
		return err
	}
	content := string(bytes)
	m := &model{
		content:  content,
		args:     make(map[string]string),
		includes: make([]string, 0),
		excludes: make([]string, 0),
		entitys:  make([]entity, 0),
	}

	if err := m.init(); err != nil {
		log.Println(fmt.Sprintf("model init error: %s", logger.ErrorWW(err)))
		return err
	}
	err = m.build(a.Output)
	if err != nil {
		log.Println(fmt.Sprintf("model build error: %s", logger.ErrorWW(err)))
		return err
	}

	return nil
}

type model struct {
	content string

	args     map[string]string
	includes []string
	excludes []string

	entitys []entity
}

type entity struct {
	name   string
	desc   string
	fields []field

	primary []string
	udxmap  map[string][]string
	idxmap  map[string][]string
	fkmap   map[string][]foreign
}

type field struct {
	name string // 字段
	desc string // 说明
	typx string // 类型
	memo string // 备注
	sql  string // sql文,只取第一个","号前的内容
}

type foreign struct {
	column1 string
	table   string
	column2 string
}

func (a *model) init() error {
	reg := regexp.MustCompile("```sql model[\\s\\S]+?```")
	ct1 := reg.FindString(a.content)
	if ct1 != "" {
		ct1 = ct1[11 : len(ct1)-3]
		reg1 := regexp.MustCompile("ARG \\w+=\\w+")
		ct2 := reg1.FindAllString(ct1, -1)
		for _, ct3 := range ct2 {
			ct3 = ct3[4:]
			o := strings.Index(ct3, "=")
			key := strings.TrimSpace(ct3[:o])
			val := strings.TrimSpace(ct3[o+1:])
			a.args[key] = val
		}

		reg2 := regexp.MustCompile("includes=[\\w, ]+")
		ct3 := reg2.FindString(ct1)
		if ct3 != "" {
			ct3 = ct3[9:]
			a.includes = strings.Split(ct3, ",")
		}
		reg3 := regexp.MustCompile("excludes=[\\w, ]+")
		ct4 := reg3.FindString(ct1)
		if ct4 != "" {
			ct4 = ct4[9:]
			a.excludes = strings.Split(ct4, ",")
		}
	}
	reg1 := regexp.MustCompile("## [\\s\\S]+?\\n---")
	ct2 := reg1.FindAllString(a.content, -1)
	for _, ct3 := range ct2 {
		ct3 = ct3[3 : len(ct3)-4]
		lines := strings.Split(ct3, "\n")

		e := entity{
			fields:  make([]field, 0),
			primary: make([]string, 0),
			udxmap:  make(map[string][]string),
			idxmap:  make(map[string][]string),
			fkmap:   make(map[string][]foreign),
		}
		for idx, ct4 := range lines {
			if idx == 0 {
				l := strings.Index(ct4, "`")
				r := strings.LastIndex(ct4, "`")
				if l < 0 || r < 0 {
					log.Println("内容无法处理:" + ct4)
					break
				}
				e.name = ct4[l+1 : r]
				e.desc = ct4[:l-1]
			} else if strings.HasPrefix(ct4, "| -") {
				// 注释跳过
				continue
			} else if strings.HasPrefix(ct4, "| 字段") {
				// 表头
				continue
			} else if strings.HasPrefix(ct4, "| ") {
				reg2 := regexp.MustCompile("\\|[^\\|]+")
				offset := reg2.FindAllStringSubmatchIndex(ct4, -1)
				if len(offset) != 5 {
					continue
				}

				f := field{}
				f.name = strings.TrimSpace(ct4[offset[0][0]+1 : offset[0][1]])
				f.desc = strings.TrimSpace(ct4[offset[1][0]+1 : offset[1][1]])
				f.typx = strings.TrimSpace(ct4[offset[2][0]+1 : offset[2][1]])
				f.memo = strings.TrimSpace(ct4[offset[3][0]+1 : offset[3][1]])

				sql := strings.TrimSpace(ct4[offset[4][0]+1 : offset[4][1]])

				sql1 := strings.Split(sql, ",")
				if len(sql1) > 1 {
					sql = strings.TrimSpace(sql1[0])
					for sidx, sql2 := range sql1 {
						if sidx == 0 {
							continue
						}
						sql2 = strings.TrimSpace(sql2)
						if strings.HasPrefix(sql2, "udx_") {
							e.udxmap[sql2] = append(e.udxmap[sql2], f.name)
						} else if strings.HasPrefix(sql2, "idx_") {
							e.idxmap[sql2] = append(e.idxmap[sql2], f.name)
						} else if strings.HasPrefix(sql2, "fk_") {
							fk := foreign{}
							fk.column1 = f.name
							sql3 := strings.Split(sql2, "->")
							sql4 := strings.Split(sql3[1], ".")
							sql5 := sql3[0]
							fk.table = sql4[0]
							fk.column2 = sql4[1]
							e.fkmap[sql5] = append(e.fkmap[sql5], fk)
						} else if sql2 == "primary" {
							e.primary = append(e.primary, f.name)
						}
					}
				}
				if !strings.Contains(strings.ToUpper(sql), "NULL") &&
					!strings.Contains(strings.ToUpper(sql), "DEFAULT") {
					sql = sql + " DEFAULT NULL"
				}
				f.sql = sql

				e.fields = append(e.fields, f)
			}
		}
		if e.name != "" && len(e.fields) > 0 {
			a.entitys = append(a.entitys, e)
		}

	}
	return nil
}

func (a *model) build(output string) error {
	for _, ct1 := range a.entitys {
		if len(a.includes) > 0 && !hasKey(ct1.name, a.includes) {
			continue
		}
		err := ct1.build(output)
		if err != nil {
			return err
		}
	}

	return nil
}

func hasKey(key string, keys []string) bool {
	for _, k := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func fixName(name string) string {
	ns := strings.Split(name, "_")
	re := ""
	for _, n := range ns {
		re += strings.ToUpper(n[:1]) + n[1:]
	}
	return re
}

func (a *entity) build(output string) error {
	content := strings.ReplaceAll(modelx, "${name}", fixName(a.name))
	content = strings.ReplaceAll(content, "${code-config}", "		Table: \""+a.name+"\",")
	// content = strings.ReplaceAll(content, "${code-edges}", "")

	cfs := ""
	for _, f := range a.fields {
		if f.name == "id" {
			continue
		}
		cfs += "		field."
		switch f.typx {
		case "数值":
			cfs += "Int"
		case "字符串":
			cfs += "String"
		case "时间格式":
			content = strings.ReplaceAll(content, "${time}", "\"time\"")
			cfs += "Time"
		}
		cfs += "(\"" + f.name + "\")"
		switch f.name {
		case "created_at", "updated_at":
			cfs += ".Default(time.Now)"
		case "version":
			cfs += ".Default(1)"
		}
		cfs += ", // " + f.desc + "\n"
	}
	if cfs != "" {
		cfs = cfs[:len(cfs)-1]
	}

	content = strings.ReplaceAll(content, "${time}", "")
	content = strings.ReplaceAll(content, "${code-fields}", cfs)

	filename := strings.ReplaceAll(a.name, "_", "")
	file, err := os.Create(output + "/" + filename + ".go")
	if err != nil {
		log.Println(fmt.Sprintf("model output error: %s", logger.ErrorWW(err)))
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		log.Println(fmt.Sprintf("model write error: %s", logger.ErrorWW(err)))
		return err
	}
	return nil
}

//  ### config ${code-config}
//  		Table: "demo"
//  ### fields ${code-fields}
//  		field.Int("status").Min(1).Max(2).Default(1),
//  ### edges ${code-edges}
//  		edge.To("children", Demo.Type),
var modelx = `
package schema

/*
  ${name}
  ID该字段内置于架构中，不需要声明。
  在基于 SQL 的数据库中，其类型默认为数据库中自动递增
  https://entgo.io/docs/getting-started/
*/
import (
	${time}

	"github.com/facebookincubator/ent"
	//"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// ${name} holds the schema definition for the ${name} entity.
type ${name} struct {
	ent.Schema
}

// Config of the ${name}.
func (${name}) Config() ent.Config {
	return ent.Config{
${code-config}
	}
}

// Hooks of the Card.
func (${name}) Hooks() []ent.Hook {
	return nil
}

// Fields of the ${name}.
func (${name}) Fields() []ent.Field {
	return []ent.Field{
${code-fields}
	}
}

// Edges of the ${name}.
func (${name}) Edges() []ent.Edge {
//	return []ent.Edge{
//	}
	return nil
}
`
