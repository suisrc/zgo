package schema

/*
  User
  ID该字段内置于架构中，不需要声明。
  在基于 SQL 的数据库中，其类型默认为数据库中自动递增
  https://entgo.io/docs/getting-started/
*/
import (
	"github.com/facebookincubator/ent"
	//"github.com/facebookincubator/ent/schema/edge"
	"github.com/facebookincubator/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Config of the User.
func (User) Config() ent.Config {
	return ent.Config{
		Table: "user",
	}
}

// Hooks of the Card.
func (User) Hooks() []ent.Hook {
	return nil
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("uid"),  // 唯一标识
		field.String("name"), // 用户名
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
