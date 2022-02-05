package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/google/uuid"
)

// LoginHistory holds the schema definition for the LoginHistory entity.
type LoginHistory struct {
	ent.Schema
}

// Fields of the LoginHistory.
func (LoginHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}),
		field.String("client_ip"),
		field.String("user_agent"),
		field.String("location").Optional(),
		field.Uint32("create_at").
			DefaultFunc(func() uint32 {
				return uint32(time.Now().Unix())
			}),
	}
}

// Edges of the LoginHistory.
func (LoginHistory) Edges() []ent.Edge {
	return nil
}
