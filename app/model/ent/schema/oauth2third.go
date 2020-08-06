package schema

/*
  Oauth2Third
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

// Oauth2Third holds the schema definition for the Oauth2Third entity.
type Oauth2Third struct {
	ent.Schema
}

// Config of the Oauth2Third.
func (Oauth2Third) Config() ent.Config {
	return ent.Config{
		Table: "oauth2_third",
	}
}

// Hooks of the Card.
func (Oauth2Third) Hooks() []ent.Hook {
	return nil
}

// Fields of the Oauth2Third.
func (Oauth2Third) Fields() []ent.Field {
	return []ent.Field{
		field.String("platform"),                   // 平台
		field.String("agent_id"),                   // 代理商标识
		field.String("suite_id"),                   // 套件标识
		field.String("app_id"),                     // 应用标识
		field.String("app_secret"),                 // 应用密钥
		field.String("authorize_url"),              // 认证地址
		field.String("token_url"),                  // 令牌地址
		field.String("profile_url"),                // 个人资料地址
		field.String("domain_def"),                 // 默认域名
		field.String("domain_check"),               // 域名认证
		field.String("js_secret"),                  // javascript密钥
		field.String("state_secret"),               // 回调state密钥
		field.Int("callback"),                      // 是否支持回调
		field.Int("cb_encrypt"),                    // 回调是否加密
		field.String("cb_token"),                   // 加密令牌
		field.String("cb_encoding"),                // 加密编码
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

// Edges of the Oauth2Third.
func (Oauth2Third) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
