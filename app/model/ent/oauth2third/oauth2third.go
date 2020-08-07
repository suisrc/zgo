// Code generated by entc, DO NOT EDIT.

package oauth2third

import (
	"time"
)

const (
	// Label holds the string label denoting the oauth2third type in the database.
	Label = "oauth2third"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldPlatform holds the string denoting the platform field in the database.
	FieldPlatform = "platform"
	// FieldAgentID holds the string denoting the agent_id field in the database.
	FieldAgentID = "agent_id"
	// FieldSuiteID holds the string denoting the suite_id field in the database.
	FieldSuiteID = "suite_id"
	// FieldAppID holds the string denoting the app_id field in the database.
	FieldAppID = "app_id"
	// FieldAppSecret holds the string denoting the app_secret field in the database.
	FieldAppSecret = "app_secret"
	// FieldAuthorizeURL holds the string denoting the authorize_url field in the database.
	FieldAuthorizeURL = "authorize_url"
	// FieldTokenURL holds the string denoting the token_url field in the database.
	FieldTokenURL = "token_url"
	// FieldProfileURL holds the string denoting the profile_url field in the database.
	FieldProfileURL = "profile_url"
	// FieldDomainDef holds the string denoting the domain_def field in the database.
	FieldDomainDef = "domain_def"
	// FieldDomainCheck holds the string denoting the domain_check field in the database.
	FieldDomainCheck = "domain_check"
	// FieldJsSecret holds the string denoting the js_secret field in the database.
	FieldJsSecret = "js_secret"
	// FieldStateSecret holds the string denoting the state_secret field in the database.
	FieldStateSecret = "state_secret"
	// FieldCallback holds the string denoting the callback field in the database.
	FieldCallback = "callback"
	// FieldCbEncrypt holds the string denoting the cb_encrypt field in the database.
	FieldCbEncrypt = "cb_encrypt"
	// FieldCbToken holds the string denoting the cb_token field in the database.
	FieldCbToken = "cb_token"
	// FieldCbEncoding holds the string denoting the cb_encoding field in the database.
	FieldCbEncoding = "cb_encoding"
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

	// Table holds the table name of the oauth2third in the database.
	Table = "oauth2_third"
)

// Columns holds all SQL columns for oauth2third fields.
var Columns = []string{
	FieldID,
	FieldPlatform,
	FieldAgentID,
	FieldSuiteID,
	FieldAppID,
	FieldAppSecret,
	FieldAuthorizeURL,
	FieldTokenURL,
	FieldProfileURL,
	FieldDomainDef,
	FieldDomainCheck,
	FieldJsSecret,
	FieldStateSecret,
	FieldCallback,
	FieldCbEncrypt,
	FieldCbToken,
	FieldCbEncoding,
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