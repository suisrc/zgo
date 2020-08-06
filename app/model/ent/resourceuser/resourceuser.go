// Code generated by entc, DO NOT EDIT.

package resourceuser

import (
	"time"
)

const (
	// Label holds the string label denoting the resourceuser type in the database.
	Label = "resource_user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldResource holds the string denoting the resource field in the database.
	FieldResource = "resource"
	// FieldCreator holds the string denoting the creator field in the database.
	FieldCreator = "creator"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"

	// Table holds the table name of the resourceuser in the database.
	Table = "resource_user"
)

// Columns holds all SQL columns for resourceuser fields.
var Columns = []string{
	FieldID,
	FieldUserID,
	FieldResource,
	FieldCreator,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldVersion,
}

var (
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
	// DefaultVersion holds the default value on creation for the version field.
	DefaultVersion int
)
