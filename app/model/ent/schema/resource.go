package schema

/*
  Resource
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

// Resource holds the schema definition for the Resource entity.
type Resource struct {
	ent.Schema
}

// Config of the Resource.
func (Resource) Config() ent.Config {
	return ent.Config{
		Table: "resource",
	}
}

// Hooks of the Card.
func (Resource) Hooks() []ent.Hook {
	return nil
}

// Fields of the Resource.
func (Resource) Fields() []ent.Field {
	return []ent.Field{
		field.String("resource"),                   // 资源名
		field.String("path"),                       // 路径
		field.String("netmask"),                    // 网络掩码
		field.Int("allow"),                         // 允许vs拒绝
		field.String("desc"),                       // 描述
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
	}
}

// Edges of the Resource.
func (Resource) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
