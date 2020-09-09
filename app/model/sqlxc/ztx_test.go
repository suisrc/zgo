package sqlxc

import (
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNowTime(t *testing.T) {
	var x sql.NullTime
	obj := NewNowTime(reflect.TypeOf(x))

	log.Println(obj)

	assert.NotNil(t, nil)
}

func TestStruct2Map(t *testing.T) {

	o2a := SigninGpaOAuth2Account{
		AccountID: 1,
		ClientID:  sql.NullInt64{Valid: false},
		ClientKID: sql.NullString{Valid: true, String: "2"},
		UserKID:   "7",
		RoleKID:   sql.NullString{Valid: true, String: "3"},
		Expired:   sql.NullInt64{Valid: true, Int64: 3},
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
		Secret:    sql.NullString{Valid: true, String: "6"},
		Status:    true,
	}

	m := Struct2Map(o2a)
	log.Println(m)

	assert.Nil(t, nil)
}

func TestSelectColumns(t *testing.T) {

	o2a := SigninGpaOAuth2Account{}

	m := SelectColumns(&o2a, "t.")
	log.Println(m)

	assert.Nil(t, nil)
}
func TestCreateUpdateSQLByNamed(t *testing.T) {

	o2a := SigninGpaOAuth2Account{
		AccountID: 1,
		ClientID:  sql.NullInt64{Valid: false},
		ClientKID: sql.NullString{Valid: true, String: "2"},
		UserKID:   "7",
		RoleKID:   sql.NullString{Valid: true, String: "3"},
		Expired:   sql.NullInt64{Valid: true, Int64: 3},
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
		Secret:    sql.NullString{Valid: true, String: "6"},
		Status:    true,
	}

	sql, pas, err := CreateUpdateSQLByNamed("table", "id", IDC{ID: 1}, o2a, nil)

	log.Println(sql)
	log.Println(pas)
	assert.Nil(t, err)
}

func TestCreateUpdateSQLByNamedSkipNil(t *testing.T) {

	o2a := SigninGpaOAuth2Account{
		AccountID: 1,
		ClientID:  sql.NullInt64{Valid: false},
		ClientKID: sql.NullString{Valid: true, String: "2"},
		UserKID:   "7",
		RoleKID:   sql.NullString{Valid: true, String: "3"},
		Expired:   sql.NullInt64{Valid: true, Int64: 3},
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
		Secret:    sql.NullString{Valid: true, String: "6"},
		Status:    true,
	}

	sql, pas, err := CreateUpdateSQLByNamedAndSkipNil("table", "id", IDC{ID: 1}, &o2a)

	log.Println(sql)
	log.Println(pas)
	assert.Nil(t, err)
}

// SigninGpaOAuth2Account account
type SigninGpaOAuth2Account struct {
	ID        int            `db:"id"`
	AccountID int            `db:"account_id"`
	ClientID  sql.NullInt64  `db:"client_id"`
	ClientKID sql.NullString `db:"client_kid"`
	UserKID   string         `db:"user_kid"`
	RoleKID   sql.NullString `db:"role_kid"`
	Expired   sql.NullInt64  `db:"expired"`
	LastIP    sql.NullString `db:"last_ip"`
	LastAt    sql.NullTime   `db:"last_at"`
	LimitExp  sql.NullTime   `db:"limit_exp"`
	LimitKey  sql.NullString `db:"limit_key"`
	Mode      sql.NullString `db:"mode"`
	Secret    sql.NullString `db:"secret"`
	Status    bool           `db:"-"`
}
