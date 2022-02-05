// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// LoginHistoriesColumns holds the columns for the "login_histories" table.
	LoginHistoriesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "app_id", Type: field.TypeUUID},
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "client_ip", Type: field.TypeString},
		{Name: "user_agent", Type: field.TypeString},
		{Name: "location", Type: field.TypeString, Nullable: true},
		{Name: "create_at", Type: field.TypeUint32},
	}
	// LoginHistoriesTable holds the schema information for the "login_histories" table.
	LoginHistoriesTable = &schema.Table{
		Name:       "login_histories",
		Columns:    LoginHistoriesColumns,
		PrimaryKey: []*schema.Column{LoginHistoriesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		LoginHistoriesTable,
	}
)

func init() {
}
