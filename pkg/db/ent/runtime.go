// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/NpoolPlatform/login-gateway/pkg/db/ent/loginhistory"
	"github.com/NpoolPlatform/login-gateway/pkg/db/ent/schema"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	loginhistoryFields := schema.LoginHistory{}.Fields()
	_ = loginhistoryFields
	// loginhistoryDescCreateAt is the schema descriptor for create_at field.
	loginhistoryDescCreateAt := loginhistoryFields[6].Descriptor()
	// loginhistory.DefaultCreateAt holds the default value on creation for the create_at field.
	loginhistory.DefaultCreateAt = loginhistoryDescCreateAt.Default.(func() uint32)
	// loginhistoryDescID is the schema descriptor for id field.
	loginhistoryDescID := loginhistoryFields[0].Descriptor()
	// loginhistory.DefaultID holds the default value on creation for the id field.
	loginhistory.DefaultID = loginhistoryDescID.Default.(func() uuid.UUID)
}
