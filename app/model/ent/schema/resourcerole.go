
package schema

/*
  ResourceRole
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

// ResourceRole holds the schema definition for the ResourceRole entity.
type ResourceRole struct {
	ent.Schema
}

// Config of the ResourceRole.
func (ResourceRole) Config() ent.Config {
	return ent.Config{
		Table: "resource_role",
	}
}

// Hooks of the Card.
func (ResourceRole) Hooks() []ent.Hook {
	return nil
}

// Fields of the ResourceRole.
func (ResourceRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"), // 角色
		field.String("resource"), // 资源名
		field.String("creator"), // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1), // 数据版本
	}
}

// Edges of the ResourceRole.
func (ResourceRole) Edges() []ent.Edge {
//	return []ent.Edge{
//	}
	return nil
}
