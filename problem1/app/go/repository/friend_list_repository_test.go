package repository

import (
	"database/sql"
	"testing"

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

func newFriendList() *model.FriendList {
	return &model.FriendList{
		Friends: []*model.User{
			{
				Id:   111111,
				Name: "hoge",
			},
			{
				Id:   222222,
				Name: "fuga",
			},
		},
	}
}

type testUser struct {
	userId int
	name   string
}

func newTestUsers() []testUser {
	return []testUser{
		{
			userId: testutil.UserIDForDebug,
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
		{
			userId: 333333,
			name:   "bar",
		},
	}
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

type userLink struct {
	user1Id int
	user2Id int
}

func newTestUserLink() []userLink {
	return []userLink{
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 111111,
		},
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 222222,
		},
		{
			user1Id: testutil.UserIDForDebug,
			user2Id: 333333,
		},
	}

}

func (r *friendListRepositoryTest) insertTestFriendLink(t *testing.T, db *sql.DB, ul userLink) {
	t.Helper()

	const q = `
	INSERT INTO friend_link (id, user1_id, user2_id)
	VALUES (0, ?, ?)`

	testRecord := []any{
		ul.user1Id,
		ul.user2Id,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

func (r *friendListRepositoryTest) insertTestBlockList(t *testing.T, db *sql.DB, ul userLink) {
	t.Helper()

	const q = `
	INSERT INTO block_list (id, user1_id, user2_id)
	VALUES (0, ?, ?)`

	testRecord := []any{
		ul.user1Id,
		ul.user2Id,
	}
	testutil.ValidateSQLArgs(t, q, testRecord...)
	testutil.ExecSQL(t, db, q, testRecord...)
}

func Test_friendListRepository_CheckUserExist(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		userId  int
		want    bool
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				tu := testUser{
					userId: testutil.UserIDForDebug,
					name:   testutil.UserNameForDebug,
				}
				rt.insertTestUserList(t, rt.db, tu)
			},
			userId:  testutil.UserIDForDebug,
			want:    true,
			wantErr: false,
		},
		{
			name:    "ok: user not exist",
			prepare: func(rt *friendListRepositoryTest) {},
			userId:  testutil.UserIDForDebug,
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.CheckUserExist(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetOneHopFriendsUserIdList(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		userId  int
		want    []int
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:  testutil.UserIDForDebug,
			want:    []int{111111, 222222, 333333},
			wantErr: false,
		},
		{
			name: "ok: no 1hop friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			userId:  testutil.UserIDForDebug,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetOneHopFriendsUserIdList(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetOneHopFrinedsUserIdList() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetBlockUsersIdList(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		userId  int
		want    []int
		wantErr bool
	}{
		{
			name: "ok: 1 user blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
				rt.insertTestBlockList(t, rt.db, userLink{
					user1Id: testutil.UserIDForDebug,
					user2Id: 111111,
				})
			},
			userId:  testutil.UserIDForDebug,
			want:    []int{111111},
			wantErr: false,
		},
		{
			name: "ok: all users blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
				for _, ul := range testUserLink {
					rt.insertTestBlockList(t, rt.db, ul)
				}
			},
			userId:  testutil.UserIDForDebug,
			want:    []int{111111, 222222, 333333},
			wantErr: false,
		},
		{
			name: "ok: no user blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:  testutil.UserIDForDebug,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetBlockUsersIdList(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBlockUsersIdList() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListByUserId(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name    string
		prepare func(*friendListRepositoryTest)
		userId  int
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for i := 0; i < 2; i++ {
					rt.insertTestFriendLink(t, rt.db, testUserLink[i])
				}
			},
			userId:  testutil.UserIDForDebug,
			want:    newFriendList(),
			wantErr: false,
		},
		{
			name: "ok: have no friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			userId: testutil.UserIDForDebug,
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListByUserId(tt.userId)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListByUserIdExcludingBlockUsers(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := newTestUserLink()

	tests := []struct {
		name       string
		prepare    func(*friendListRepositoryTest)
		userId     int
		blockUsers []int
		want       *model.FriendList
		wantErr    bool
	}{
		{
			name: "ok: 1 friend blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:     testutil.UserIDForDebug,
			blockUsers: []int{333333},
			want:       newFriendList(),
			wantErr:    false,
		},
		{
			name: "ok: all friends blocked",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:     testutil.UserIDForDebug,
			blockUsers: []int{111111, 222222, 333333},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name:       "ng: blockUsers nil",
			prepare:    func(rt *friendListRepositoryTest) {},
			userId:     testutil.UserIDForDebug,
			blockUsers: nil,
			want:       nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListByUserIdExcludingBlockUsers(tt.userId, tt.blockUsers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserIdExcludingBlockUsers() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListOfFriendsByUserId(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := []userLink{
		{
			user1Id: 123456789,
			user2Id: 444444,
		},
		{
			user1Id: 444444,
			user2Id: 111111,
		},
		{
			user1Id: 444444,
			user2Id: 222222,
		},
		{
			user1Id: 444444,
			user2Id: 333333,
		},
	}
	testUserLink2 := newTestUserLink()

	tests := []struct {
		name         string
		prepare      func(*friendListRepositoryTest)
		userId       int
		excludeUsers []int
		want         *model.FriendList
		wantErr      bool
	}{
		{
			name: "ok: 1 friend excluded",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{333333},
			want:         newFriendList(),
			wantErr:      false,
		},
		{
			name: "ok: all friends excluded",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{111111, 222222, 333333},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name: "ok: have no 2hop friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				for _, ul := range testUserLink2 {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{444444},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name: "ok: have no friend",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{111111},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name:         "ng: excludeUsers nil",
			prepare:      func(rt *friendListRepositoryTest) {},
			userId:       testutil.UserIDForDebug,
			excludeUsers: nil,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListOfFriendsByUserId(tt.userId, tt.excludeUsers)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListRepository_GetFriendListOfFriendsByUserIdWithPaging(t *testing.T) {
	testUsers := newTestUsers()
	testUserLink := []userLink{
		{
			user1Id: 123456789,
			user2Id: 444444,
		},
		{
			user1Id: 444444,
			user2Id: 111111,
		},
		{
			user1Id: 444444,
			user2Id: 222222,
		},
		{
			user1Id: 444444,
			user2Id: 333333,
		},
	}

	tests := []struct {
		name         string
		prepare      func(*friendListRepositoryTest)
		userId       int
		excludeUsers []int
		limit        int
		offset       int
		want         *model.FriendList
		wantErr      bool
	}{
		{
			name: "ok: limit",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{444444},
			limit:        2,
			offset:       0,
			want:         newFriendList(),
			wantErr:      false,
		},
		{
			name: "ok: offset",
			prepare: func(rt *friendListRepositoryTest) {
				for _, tu := range testUsers {
					rt.insertTestUserList(t, rt.db, tu)
				}
				rt.insertTestUserList(t, rt.db, testUser{
					userId: 444444,
					name:   "piyo",
				})
				for _, ul := range testUserLink {
					rt.insertTestFriendLink(t, rt.db, ul)
				}
			},
			userId:       testutil.UserIDForDebug,
			excludeUsers: []int{444444},
			limit:        3,
			offset:       2,
			want: &model.FriendList{
				Friends: []*model.User{
					{
						Id:   333333,
						Name: "bar",
					},
				},
			},
			wantErr: false,
		},
		{
			name:         "ng: excludeUsers nil",
			prepare:      func(rt *friendListRepositoryTest) {},
			userId:       testutil.UserIDForDebug,
			excludeUsers: nil,
			want:         nil,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := newFriendListRepositoryTest(t)
			tt.prepare(rt)

			got, err := rt.flr.GetFriendListOfFriendsByUserIdWithPaging(tt.userId, tt.excludeUsers, tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserIdWithPaging() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
