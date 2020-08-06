// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/suisrc/zgo/app/model/ent/userdetail"
)

// UserDetail is the model entity for the UserDetail schema.
type UserDetail struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// UserID holds the value of the "user_id" field.
	UserID int `json:"user_id,omitempty"`
	// Nickname holds the value of the "nickname" field.
	Nickname string `json:"nickname,omitempty"`
	// Avatar holds the value of the "avatar" field.
	Avatar string `json:"avatar,omitempty"`
	// Creator holds the value of the "creator" field.
	Creator string `json:"creator,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Version holds the value of the "version" field.
	Version int `json:"version,omitempty"`
	// String1 holds the value of the "string_1" field.
	String1 string `json:"string_1,omitempty"`
	// String2 holds the value of the "string_2" field.
	String2 string `json:"string_2,omitempty"`
	// String3 holds the value of the "string_3" field.
	String3 string `json:"string_3,omitempty"`
	// Number1 holds the value of the "number_1" field.
	Number1 int `json:"number_1,omitempty"`
	// Number2 holds the value of the "number_2" field.
	Number2 int `json:"number_2,omitempty"`
	// Number3 holds the value of the "number_3" field.
	Number3 int `json:"number_3,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserDetail) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullInt64{},  // user_id
		&sql.NullString{}, // nickname
		&sql.NullString{}, // avatar
		&sql.NullString{}, // creator
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullInt64{},  // version
		&sql.NullString{}, // string_1
		&sql.NullString{}, // string_2
		&sql.NullString{}, // string_3
		&sql.NullInt64{},  // number_1
		&sql.NullInt64{},  // number_2
		&sql.NullInt64{},  // number_3
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserDetail fields.
func (ud *UserDetail) assignValues(values ...interface{}) error {
	if m, n := len(values), len(userdetail.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ud.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field user_id", values[0])
	} else if value.Valid {
		ud.UserID = int(value.Int64)
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field nickname", values[1])
	} else if value.Valid {
		ud.Nickname = value.String
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field avatar", values[2])
	} else if value.Valid {
		ud.Avatar = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field creator", values[3])
	} else if value.Valid {
		ud.Creator = value.String
	}
	if value, ok := values[4].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[4])
	} else if value.Valid {
		ud.CreatedAt = value.Time
	}
	if value, ok := values[5].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[5])
	} else if value.Valid {
		ud.UpdatedAt = value.Time
	}
	if value, ok := values[6].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field version", values[6])
	} else if value.Valid {
		ud.Version = int(value.Int64)
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_1", values[7])
	} else if value.Valid {
		ud.String1 = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_2", values[8])
	} else if value.Valid {
		ud.String2 = value.String
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field string_3", values[9])
	} else if value.Valid {
		ud.String3 = value.String
	}
	if value, ok := values[10].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field number_1", values[10])
	} else if value.Valid {
		ud.Number1 = int(value.Int64)
	}
	if value, ok := values[11].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field number_2", values[11])
	} else if value.Valid {
		ud.Number2 = int(value.Int64)
	}
	if value, ok := values[12].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field number_3", values[12])
	} else if value.Valid {
		ud.Number3 = int(value.Int64)
	}
	return nil
}

// Update returns a builder for updating this UserDetail.
// Note that, you need to call UserDetail.Unwrap() before calling this method, if this UserDetail
// was returned from a transaction, and the transaction was committed or rolled back.
func (ud *UserDetail) Update() *UserDetailUpdateOne {
	return (&UserDetailClient{config: ud.config}).UpdateOne(ud)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ud *UserDetail) Unwrap() *UserDetail {
	tx, ok := ud.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserDetail is not a transactional entity")
	}
	ud.config.driver = tx.drv
	return ud
}

// String implements the fmt.Stringer.
func (ud *UserDetail) String() string {
	var builder strings.Builder
	builder.WriteString("UserDetail(")
	builder.WriteString(fmt.Sprintf("id=%v", ud.ID))
	builder.WriteString(", user_id=")
	builder.WriteString(fmt.Sprintf("%v", ud.UserID))
	builder.WriteString(", nickname=")
	builder.WriteString(ud.Nickname)
	builder.WriteString(", avatar=")
	builder.WriteString(ud.Avatar)
	builder.WriteString(", creator=")
	builder.WriteString(ud.Creator)
	builder.WriteString(", created_at=")
	builder.WriteString(ud.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(ud.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", version=")
	builder.WriteString(fmt.Sprintf("%v", ud.Version))
	builder.WriteString(", string_1=")
	builder.WriteString(ud.String1)
	builder.WriteString(", string_2=")
	builder.WriteString(ud.String2)
	builder.WriteString(", string_3=")
	builder.WriteString(ud.String3)
	builder.WriteString(", number_1=")
	builder.WriteString(fmt.Sprintf("%v", ud.Number1))
	builder.WriteString(", number_2=")
	builder.WriteString(fmt.Sprintf("%v", ud.Number2))
	builder.WriteString(", number_3=")
	builder.WriteString(fmt.Sprintf("%v", ud.Number3))
	builder.WriteByte(')')
	return builder.String()
}

// UserDetails is a parsable slice of UserDetail.
type UserDetails []*UserDetail

func (ud UserDetails) config(cfg config) {
	for _i := range ud {
		ud[_i].config = cfg
	}
}
