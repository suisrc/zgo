package schema

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/zgo/app/model/sqlxc"
)

func TestCreateUpdateSQLByNamedSkipNilAndSet(t *testing.T) {

	token := AccountOAuth2Token{
		TokenID:     "1234",
		AccessToken: sql.NullString{Valid: true, String: "123"},
	}

	IDX := sqlxc.IdxColumn{Column: "token_kid", KID: token.TokenID, Update: true}
	sql, pas, err := sqlxc.CreateUpdateSQLByNamedAndSkipNilAndSet("table", IDX, &token)

	log.Println(sql)
	log.Println(pas)
	assert.Nil(t, err)
	assert.NotNil(t, nil)
}
