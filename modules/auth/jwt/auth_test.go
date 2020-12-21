package jwt

import (
	"context"
	"log"
	"testing"

	"github.com/suisrc/zgo/modules/logger"
	"github.com/suisrc/zgo/modules/store/buntdb"

	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	store, err := buntdb.NewStore(":memory:")
	assert.Nil(t, err)

	token := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZWYuY29tIiwiZXhwIjoxNTk5NTU0ODk0LCJqdGkiOiJfa3hpcHJ3djZ1eGYzaHVyemxuM3E0MjBxdjFmbGVoYTkiLCJpYXQiOjE1OTk1NTQ4ODksImlzcyI6ImFiYy5jb20iLCJuYmYiOjE1OTk1NTQ4ODksInN1YiI6IjEyMyIsIm5hbSI6Ikpzb24iLCJyb2wiOiI3ODkifQ.085VKYbb7yV5cOQZoIYMhvXyY7gChE6A_w_PoOZnVTQtzkeEnRjU_vSDazkhyZPbA3opOd9wZrpI_pRzbPw2yA"
	secret := "zgo_c1wcx9vp97pd4iesw68o3byh71fleha9"

	jwtAuth := New(store, Option(func(o *options) {
		o.signingSecret = []byte(secret)
	}))
	defer jwtAuth.Release()
	ctx := context.Background()

	r, u, err := jwtAuth.RefreshToken(ctx, token, nil)
	log.Println(r)
	log.Println(u)

	assert.Nil(t, err)
	assert.NotNil(t, nil)
}

func TestAuth(t *testing.T) {
	store, err := buntdb.NewStore(":memory:")
	assert.Nil(t, err)
	// var store Storer

	var ref TokenRef

	jwtAuth := New(store, Option(func(o *options) {
		o.tokenFunc = func(ctx context.Context) (string, error) {
			return ref.ref, nil
		}
		o.expired = 5
	}))

	defer jwtAuth.Release()

	ctx := context.Background()

	user := &UserInfo{
		UserName: "Json",
		UserID:   "123",
		RoleID:   "789",
		//TokenID:  "456",
		Issuer:   "abc.com",
		Audience: "def.com",
	}

	token, _, err := jwtAuth.GenerateToken(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	data, err := token.EncodeToJSON()
	logger.Infof(ctx, "%s", string(data))

	ref.ref = token.GetAccessToken()
	uInfo, err := jwtAuth.GetUserInfo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, user.UserID, uInfo.GetUserID())

	err = jwtAuth.DestroyToken(ctx, uInfo)
	assert.Nil(t, err)

	uInfo, err = jwtAuth.GetUserInfo(ctx)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid token")
	//assert.Empty(t, id)
	assert.NotNil(t, nil)
}

type TokenRef struct {
	ref string
}

// UserInfo 用户信息声明
type UserInfo struct {
	UserName  string
	UserID    string
	RoleID    string
	TokenID   string
	Issuer    string
	Audience  string
	AccountID string
	XID       string
	TID       string
}

// GetUserName name
func (u *UserInfo) GetUserName() string {
	return u.UserName
}

// GetUserID user
func (u *UserInfo) GetUserID() string {
	return u.UserID
}

// GetRoleID role
func (u *UserInfo) GetRoleID() string {
	return u.RoleID
}

// SetRoleID role
func (u *UserInfo) SetRoleID(r string) string {
	x := u.RoleID
	u.RoleID = r
	return x
}

// GetTokenID token
func (u *UserInfo) GetTokenID() string {
	return u.TokenID
}

// GetAccountID token
func (u *UserInfo) GetAccountID() string {
	return u.AccountID
}

// GetIssuer issuer
func (u *UserInfo) GetIssuer() string {
	return u.Issuer
}

// GetAudience audience
func (u *UserInfo) GetAudience() string {
	return u.Audience
}

// GetProps props
func (u *UserInfo) GetProps() (interface{}, bool) {
	return nil, false
}

// GetXID xid
func (u *UserInfo) GetXID() string {
	return u.XID
}

// GetTID tid
func (u *UserInfo) GetTID() string {
	return u.TID
}
