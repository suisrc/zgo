
package schema

/*
  UserMessage
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

// UserMessage holds the schema definition for the UserMessage entity.
type UserMessage struct {
	ent.Schema
}

// Config of the UserMessage.
func (UserMessage) Config() ent.Config {
	return ent.Config{
		Table: "user_message",
	}
}

// Hooks of the Card.
func (UserMessage) Hooks() []ent.Hook {
	return nil
}

// Fields of the UserMessage.
func (UserMessage) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid"), // 索引
		field.String("avatar"), // 头像
		field.String("title"), // 标题
		field.String("datetime"), // 日期
		field.String("type"), // 类型
		field.Int("read"), // 已读
		field.String("description"), // 描述
		field.Int("clickClose"), // 关闭
		field.Int("status"), // 状态
		field.String("creator"), // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1), // 数据版本
	}
}

// Edges of the UserMessage.
func (UserMessage) Edges() []ent.Edge {
//	return []ent.Edge{
//	}
	return nil
}
