package schema

/*
  MenuRole
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

// MenuRole holds the schema definition for the MenuRole entity.
type MenuRole struct {
	ent.Schema
}

// Config of the MenuRole.
func (MenuRole) Config() ent.Config {
	return ent.Config{
		Table: "menu_role",
	}
}

// Hooks of the Card.
func (MenuRole) Hooks() []ent.Hook {
	return nil
}

// Fields of the MenuRole.
func (MenuRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("role_id"),                       // 角色 ID
		field.Int("user_id"),                       // 用户 ID
		field.Int("menu_id"),                       // 菜单 ID
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
	}
}

// Edges of the MenuRole.
func (MenuRole) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
