package schema

/*
  UserRole
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

// UserRole holds the schema definition for the UserRole entity.
type UserRole struct {
	ent.Schema
}

// Config of the UserRole.
func (UserRole) Config() ent.Config {
	return ent.Config{
		Table: "user_role",
	}
}

// Hooks of the Card.
func (UserRole) Hooks() []ent.Hook {
	return nil
}

// Fields of the UserRole.
func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),                       // 账户标识
		field.Int("role_id"),                       // 客户端标识
		field.Time("expired"),                      // 授权有效期
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
	}
}

// Edges of the UserRole.
func (UserRole) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
