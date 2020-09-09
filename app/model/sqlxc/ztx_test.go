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
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
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
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
	}

	sql, pas, err := CreateUpdateSQLByNamed("table", "id", IDC{ID: 1}, o2a, nil, nil)

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
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
	}

	sql, pas, err := CreateUpdateSQLByNamedAndSkipNil("table", "id", IDC{ID: 1}, &o2a)

	log.Println(sql)
	log.Println(pas)
	assert.Nil(t, err)
	assert.NotNil(t, nil)
}

func TestCreateUpdateSQLByNamedSkipNilAndSet(t *testing.T) {

	o2a := SigninGpaOAuth2Account{
		AccountID: 1,
		ClientID:  sql.NullInt64{Valid: false},
		ClientKID: sql.NullString{Valid: true, String: "2"},
		UserKID:   "7",
		RoleKID:   sql.NullString{Valid: true, String: "3"},
		LastIP:    sql.NullString{Valid: true, String: "5"},
		LastAt:    sql.NullTime{Valid: true, Time: time.Now()},
		LimitExp:  sql.NullTime{Valid: false},
		LimitKey:  sql.NullString{Valid: false},
		Mode:      sql.NullString{Valid: true, String: "signin"},
	}

	sql, pas, err := CreateUpdateSQLByNamedAndSkipNilAndSet("table", "id", IDC{ID: 1}, &o2a)

	log.Println(sql)
	log.Println(pas)
	assert.Nil(t, err)
	assert.NotNil(t, nil)
}

// SigninGpaOAuth2Account account
type SigninGpaOAuth2Account struct {
	ID           int            `db:"id"`
	AccountID    int            `db:"account_id"`
	UserKID      string         `db:"user_kid"`
	TokenID      string         `db:"token_kid"`
	ClientID     sql.NullInt64  `db:"client_id"`
	ClientKID    sql.NullString `db:"client_kid"`
	RoleKID      sql.NullString `db:"role_kid"`
	LastIP       sql.NullString `db:"last_ip"`
	LastAt       sql.NullTime   `db:"last_at"`
	LimitExp     sql.NullTime   `db:"limit_exp"`
	LimitKey     sql.NullString `db:"limit_key"`
	Mode         sql.NullString `db:"mode"`
	ExpiresAt    sql.NullInt64  `db:"expires_at"`
	AccessToken  sql.NullString `db:"access_token"`
	RefreshToken sql.NullString `db:"refresh_token"`
	RefreshCount sql.NullInt64  `db:"refresh_count" set:"=+1"`
	Status       sql.NullBool   `db:"status"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	Version      sql.NullInt64  `db:"version" set:"=+1"`
}
