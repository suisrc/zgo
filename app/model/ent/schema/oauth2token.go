
package schema

/*
  Oauth2Token
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

// Oauth2Token holds the schema definition for the Oauth2Token entity.
type Oauth2Token struct {
	ent.Schema
}

// Config of the Oauth2Token.
func (Oauth2Token) Config() ent.Config {
	return ent.Config{
		Table: "oauth2_token",
	}
}

// Hooks of the Card.
func (Oauth2Token) Hooks() []ent.Hook {
	return nil
}

// Fields of the Oauth2Token.
func (Oauth2Token) Fields() []ent.Field {
	return []ent.Field{
		field.String("oauth2_id"), // 平台
		field.String("access_token"), // 代理商标识
		field.String("expires_in"), // 有限期间隔
		field.String("create_time"), // 凭据创建时间
		field.Int("sync_lock"), // 同步锁
		field.Int("call_count"), // 调用次数
		field.String("creator"), // 创建人
		field.Time("created_at").Default(time.Now), // 创建时间
		field.Time("updated_at").Default(time.Now), // 更新时间
		field.Int("version").Default(1), // 数据版本
		field.String("string_1"), // 备用字段
		field.String("string_2"), // 备用字段
		field.String("string_3"), // 备用字段
		field.Int("number_1"), // 备用字段
		field.Int("number_2"), // 备用字段
		field.Int("number_3"), // 备用字段
	}
}

// Edges of the Oauth2Token.
func (Oauth2Token) Edges() []ent.Edge {
//	return []ent.Edge{
//	}
	return nil
}
