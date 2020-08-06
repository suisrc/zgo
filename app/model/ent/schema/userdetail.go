package schema

/*
  UserDetail
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

// UserDetail holds the schema definition for the UserDetail entity.
type UserDetail struct {
	ent.Schema
}

// Config of the UserDetail.
func (UserDetail) Config() ent.Config {
	return ent.Config{
		Table: "user_detail",
	}
}

// Hooks of the Card.
func (UserDetail) Hooks() []ent.Hook {
	return nil
}

// Fields of the UserDetail.
func (UserDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),                       // 用户 ID
		field.String("nickname"),                   // 昵称
		field.String("avatar"),                     // 头像
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

// Edges of the UserDetail.
func (UserDetail) Edges() []ent.Edge {
	//	return []ent.Edge{
	//	}
	return nil
}
