package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/shopspring/decimal"
	"net/url"
)

// Car holds the schema definition for the Car entity.
type Car struct {
	ent.Schema
}

func (Car) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Car.
func (Car) Fields() []ent.Field {
	return []ent.Field{
		field.String("model"),
		field.JSON("url", &url.URL{}).Optional(),
		field.Enum("status").Values("pending", "approved", "rejected", "cancelled").Default("pending"),
		field.Time("registered_at").Optional().Nillable().StructTag(`json:"registered_at"`),
		field.Time("approved_at").Optional().Nillable().
			// 默认是：`json:"approved_at,omitempty"`，如果是空，json序列化会不导出
			StructTag(`json:"approved_at"`),
		// 自定义金额类型
		field.Other("price", decimal.Decimal{}).
			SchemaType(map[string]string{
				dialect.MySQL: "DECIMAL(10,2)",
			}).
			Default(decimal.Decimal{}).
			Optional(),
	}
}

// Edges of the Car.
func (Car) Edges() []ent.Edge {
	return []ent.Edge{edge.From("owner", User.Type).
		Ref("cars").
		Unique(),
	}
}
