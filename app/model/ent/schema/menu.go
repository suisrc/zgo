package schema

/*
  Menu
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

// Menu holds the schema definition for the Menu entity.
type Menu struct {
	ent.Schema
}

// Config of the Menu.
func (Menu) Config() ent.Config {
	return ent.Config{
		Table: "menu",
	}
}

// Hooks of the Card.
func (Menu) Hooks() []ent.Hook {
	return nil
}

// Fields of the Menu.
func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("parent_id"),                  // 父级 ID
		field.String("name"),                       // 菜单名称
		field.Int("sequence"),                      // 排序值
		field.String("icon"),                       // 图标
		field.String("router"),                     // 访问路由
		field.String("memo"),                       // 备注
		field.Int("status"),                        // 状态
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
	}
}

// Edges of the Menu.
func (Menu) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
