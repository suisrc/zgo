// Code generated by entc, DO NOT EDIT.

package oauth2client

import (
	"time"
)

const (
	// Label holds the string label denoting the oauth2client type in the database.
	Label = "oauth2client"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldClientKey holds the string denoting the client_key field in the database.
	FieldClientKey = "client_key"
	// FieldAudience holds the string denoting the audience field in the database.
	FieldAudience = "audience"
	// FieldIssuer holds the string denoting the issuer field in the database.
	FieldIssuer = "issuer"
	// FieldExpired holds the string denoting the expired field in the database.
	FieldExpired = "expired"
	// FieldTokenType holds the string denoting the token_type field in the database.
	FieldTokenType = "token_type"
	// FieldSMethod holds the string denoting the s_method field in the database.
	FieldSMethod = "s_method"
	// FieldSSecret holds the string denoting the s_secret field in the database.
	FieldSSecret = "s_secret"
	// FieldTokenGetter holds the string denoting the token_getter field in the database.
	FieldTokenGetter = "token_getter"
	// FieldSigninURL holds the string denoting the signin_url field in the database.
	FieldSigninURL = "signin_url"
	// FieldSigninForce holds the string denoting the signin_force field in the database.
	FieldSigninForce = "signin_force"
	// FieldCreator holds the string denoting the creator field in the database.
	FieldCreator = "creator"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldVersion holds the string denoting the version field in the database.
	FieldVersion = "version"
	// FieldString1 holds the string denoting the string_1 field in the database.
	FieldString1 = "string_1"
	// FieldString2 holds the string denoting the string_2 field in the database.
	FieldString2 = "string_2"
	// FieldString3 holds the string denoting the string_3 field in the database.
	FieldString3 = "string_3"
	// FieldNumber1 holds the string denoting the number_1 field in the database.
	FieldNumber1 = "number_1"
	// FieldNumber2 holds the string denoting the number_2 field in the database.
	FieldNumber2 = "number_2"
	// FieldNumber3 holds the string denoting the number_3 field in the database.
	FieldNumber3 = "number_3"

	// Table holds the table name of the oauth2client in the database.
	Table = "oauth2_client"
)

// Columns holds all SQL columns for oauth2client fields.
var Columns = []string{
	FieldID,
	FieldClientKey,
	FieldAudience,
	FieldIssuer,
	FieldExpired,
	FieldTokenType,
	FieldSMethod,
	FieldSSecret,
	FieldTokenGetter,
	FieldSigninURL,
	FieldSigninForce,
	FieldCreator,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldVersion,
	FieldString1,
	FieldString2,
	FieldString3,
	FieldNumber1,
	FieldNumber2,
	FieldNumber3,
}

var (
	// DefaultCreatedAt holds the default value on creation for the created_at field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the updated_at field.
	DefaultUpdatedAt func() time.Time
	// DefaultVersion holds the default value on creation for the version field.
	DefaultVersion int
)