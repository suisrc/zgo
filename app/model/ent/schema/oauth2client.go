package schema

/*
  Oauth2Client
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

// Oauth2Client holds the schema definition for the Oauth2Client entity.
type Oauth2Client struct {
	ent.Schema
}

// Config of the Oauth2Client.
func (Oauth2Client) Config() ent.Config {
	return ent.Config{
		Table: "oauth2_client",
	}
}

// Hooks of the Card.
func (Oauth2Client) Hooks() []ent.Hook {
	return nil
}

// Fields of the Oauth2Client.
func (Oauth2Client) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_key"),                 // 客户端标识
		field.String("audience"),                   // 账户接受平台
		field.String("issuer"),                     // 账户签发平台
		field.Int("expired"),                       // 令牌有效期
		field.Int("token_type"),                    // 令牌类型
		field.String("s_method"),                   // 令牌方法
		field.String("s_secret"),                   // 令牌密钥
		field.String("token_getter"),               // 令牌获取方法
		field.String("signin_url"),                 // 登陆地址
		field.Int("signin_force"),                  // 强制跳转登陆
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

// Edges of the Oauth2Client.
func (Oauth2Client) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
