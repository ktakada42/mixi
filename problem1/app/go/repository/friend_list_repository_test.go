package repository

import (
	"database/sql"
	"net/http"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/model"
	"problem1/testutil"
)

type friendListRepositoryTest struct {
	db  *sql.DB
	flr FriendListRepository
}

func newFriendListRepositoryTest(t *testing.T) *friendListRepositoryTest {
	t.Helper()

	db := testutil.PrepareMySQL(t)
	flr := NewFriendListRepository(db)

	return &friendListRepositoryTest{
		db:  db,
		flr: flr,
	}
}

func newFriendList() []*model.FriendList {
	return []*model.FriendList{
		{
			Id:   111111,
			Name: "hoge",
		},
		{
			Id:   222222,
			Name: "fuga",
		},
	}
}

type testUser struct {
	userId int
	name   string
}

func (r *friendListRepositoryTest) insertTestUserList(t *testing.T, db *sql.DB, tu testUser) {
	t.Helper()

	const q = `
INSERT INTO users (id, user_id, name)
VALUES (0, ?, ?)`

	testRecord := []any{
		tu.userId,
		tu.name,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

type friendLink struct {
	user1Id int
	user2Id int
}

func (r *friendListRepositoryTest) insertTestFriendLink(t *testing.T, db *sql.DB, fl friendLink) {
	t.Helper()

	const q = `
INSERT INTO friend_link (id, user1_id, user2_id)
VALUES (0, ?, ?)`

	testRecord := []any{
		fl.user1Id,
		fl.user2Id,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

func Test_friendListRepository_GetFriendListByUserId(t *testing.T) {
	want := newFriendList()
	testUsers := []testUser{
		{
			userId: 123456789,
			name:   testutil.UserNameForDebug,
		},
		{
			userId: 111111,
			name:   "hoge",
		},
		{
			userId: 222222,
			name:   "fuga",
		},
	}
	testFriendLinks := []friendLink{
		{
			user1Id: 123456789,
			user2Id: 111111,
		},
		{
			user1Id: 123456789,
			user2Id: 222222,
		},
	}

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		param   string
		want    []*model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, fl := range testFriendLinks {
					rt.insertTestFriendLink(t, rt.db, fl)
				}
			},
			param:   "/?userId=123456789",
			want:    want,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			e := echo.New()
			e.GET("", nil)
			url, err := url.Parse(tt.param)
			if err != nil {
				t.Fatal(err)
			}
			c := e.NewContext(&http.Request{URL: url}, nil)

			got, err := rt.flr.GetFriendListByUserId(c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
