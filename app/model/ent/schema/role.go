
package schema

/*
  Role
  ID该字段内置于架构中，不需要声明。
  在基于 SQL 的数据库中，其类型默认为数据库中自动递增
  https://entgo.io/docs/getting-started/
*/
import (
	"time"

	"github.com/facebookincubator/ent"
	//"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Config of the Role.
func (Role) Config() ent.Config {
	return ent.Config{
		Table: "role",
	}
}

// Hooks of the Card.
func (Role) Hooks() []ent.Hook {
	return nil
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid"), // 唯一标识
		field.String("name"), // 角色名
		field.String("desc"), // 角色描述
		field.String("creator"), // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1), // 数据版本
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
//	return []ent.Edge{
//	}
	return nil
}
