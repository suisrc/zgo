package schema

/*
  Oauth2Account
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

// Oauth2Account holds the schema definition for the Oauth2Account entity.
type Oauth2Account struct {
	ent.Schema
}

// Config of the Oauth2Account.
func (Oauth2Account) Config() ent.Config {
	return ent.Config{
		Table: "oauth2_account",
	}
}

// Hooks of the Card.
func (Oauth2Account) Hooks() []ent.Hook {
	return nil
}

// Fields of the Oauth2Account.
func (Oauth2Account) Fields() []ent.Field {
	return []ent.Field{
		field.Int("client_id"),                     // 客户端标识
		field.String("secret"),                     // 密钥
		field.Time("expired"),                      // 授权有效期
		field.String("creator"),                    // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1),            // 数据版本
		field.String("string_1"),                   // 备用字段
		field.String("string_2"),                   // 备用字段
		field.String("string_3"),                   // 备用字段
		field.Int("number_1"),                      // 备用字段
		field.Int("number_2"),                      // 备用字段
		field.Int("number_3"),                      // 备用字段
	}
}

// Edges of the Oauth2Account.
func (Oauth2Account) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
