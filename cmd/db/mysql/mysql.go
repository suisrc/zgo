package mysql

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
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
		log.Println(fmt.Sprintf("model file error: %s", err.Error()))
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
		log.Println(fmt.Sprintf("model init error: %s", err.Error()))
		return err
	}
	content, err = m.build()
	if err != nil {
		log.Println(fmt.Sprintf("model build error: %s", err.Error()))
		return err
	}
	var file *os.File
	file, err = os.Create(a.Output)
	if err != nil {
		log.Println(fmt.Sprintf("model output error: %s", err.Error()))
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		log.Println(fmt.Sprintf("model write error: %s", err.Error()))
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
	idumap  map[string][]string
	idxmap  map[string][]string
	fkmap   map[string][]foreign
}

type field struct {
	name string // 字段
	desc string // 说明
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
			idumap:  make(map[string][]string),
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
						if strings.HasPrefix(sql2, "idu_") {
							e.idumap[sql2] = append(e.idumap[sql2], f.name)
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

func (a *model) build() (string, error) {
	content := "-- -------------------------------------------------------\n"
	content += "-- build by cmd/db/mysql/mysql.go\n-- time: " + time.Now().Format("2006-01-02 15:04:05 CST") + "\n"
	content += "-- -------------------------------------------------------\n-- 表结构"
	for _, ct1 := range a.entitys {
		sql, err := ct1.build()
		if err != nil {
			log.Println("构建结构发生异常:" + ct1.name)
			continue
		}
		content += "\n-- -------------------------------------------------------\n"
		content += sql
	}

	content += "\n-- -------------------------------------------------------\n"
	content += "\n-- -------------------------------------------------------\n-- 表外键"
	content += "\n-- -------------------------------------------------------"
	for _, ct1 := range a.entitys {
		sql, err := ct1.foreign()
		if err != nil {
			log.Println("构建外键发生异常:" + ct1.name)
			continue
		}
		if sql == "" {
			continue
		}
		content += "\n" + sql + "\n"
	}

	content += "\n-- -------------------------------------------------------"
	return content, nil
}

func (a *entity) build() (string, error) {
	content := "-- " + a.desc + "\n"
	content += "CREATE TABLE `" + a.name + "` (\n"
	for _, f := range a.fields {
		content += "  `" + f.name + "` " + f.sql
		if f.desc != "" {
			content += " COMMENT '" + f.desc + "'"
		}
		content += ",\n"
	}

	for key, val := range a.idumap {
		content += "  UNIQUE " + key + "("
		for i, p := range val {
			if i == 0 {
				content += "`" + p + "`"
			} else {
				content += ",`" + p + "`"
			}
		}
		content += "),\n"
	}

	for key, val := range a.idxmap {
		content += "  INDEX " + key + "("
		for i, p := range val {
			if i == 0 {
				content += "`" + p + "`"
			} else {
				content += ",`" + p + "`"
			}
		}
		content += "),\n"
	}

	//for key, val := range a.fkmap {
	//	content += "  FOREIGN KEY " + key + "("
	//	for i, p := range val {
	//		if i == 0 {
	//			content += "`" + p.column1 + "`"
	//		} else {
	//			content += ",`" + p.column1 + "`"
	//		}
	//	}
	//	content += "),\n"
	//	content += "  REFERENCES " + val[0].table + "("
	//	for i, p := range val {
	//		if i == 0 {
	//			content += "`" + p.column2 + "`"
	//		} else {
	//			content += ",`" + p.column2 + "`"
	//		}
	//	}
	//	content += "),\n"
	//}

	content += "  PRIMARY KEY ("
	for i, p := range a.primary {
		if i == 0 {
			content += "`" + p + "`"
		} else {
			content += ",`" + p + "`"
		}
	}
	content += ")\n"

	content += ") ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;"
	return content, nil
}

func (a *entity) foreign() (string, error) {
	if len(a.fkmap) == 0 {
		return "", nil
	}
	content := "ALTER TABLE `" + a.name + "`\n"

	for key, val := range a.fkmap {
		content += "ADD CONSTRAINT `" + key + "` FOREIGN KEY ("
		for i, p := range val {
			if i == 0 {
				content += "`" + p.column1 + "`"
			} else {
				content += ",`" + p.column1 + "`"
			}
		}
		content += ")  REFERENCES `" + val[0].table + "` ("
		for i, p := range val {
			if i == 0 {
				content += "`" + p.column2 + "`"
			} else {
				content += ",`" + p.column2 + "`"
			}
		}
		content += "),\n"
	}

	content = content[:len(content)-2] + ";"
	return content, nil
}
