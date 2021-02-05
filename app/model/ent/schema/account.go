package schema

/*
  Account
  ID该字段内置于架构中，不需要声明。
  在基于 SQL 的数据库中，其类型默认为数据库中自动递增
  https://entgo.io/docs/getting-started/
*/
import (
	"github.com/facebook/ent"
	//"github.com/facebook/ent/schema/edge"
	"github.com/facebook/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Config of the Account.
func (Account) Config() ent.Config {
	return ent.Config{
		Table: "account",
	}
}

// Hooks of the Card.
func (Account) Hooks() []ent.Hook {
	return nil
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),         // ID
		field.String("account"), // 账户
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
