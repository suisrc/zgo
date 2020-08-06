package schema

/*
  RoleRole
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

// RoleRole holds the schema definition for the RoleRole entity.
type RoleRole struct {
	ent.Schema
}

// Config of the RoleRole.
func (RoleRole) Config() ent.Config {
	return ent.Config{
		Table: "role_role",
	}
}

// Hooks of the Card.
func (RoleRole) Hooks() []ent.Hook {
	return nil
}

// Fields of the RoleRole.
func (RoleRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("owner_id"),                      // 父节点标识
		field.Int("child_id"),                      // 子节点标识
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
	}
}

// Edges of the RoleRole.
func (RoleRole) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
