package service

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/mock/mock_repository"
	"problem1/model"
	"problem1/testutil"
)

type friendListServiceTest struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
	flr  *mock_repository.MockFriendListRepository
	fls  FriendListService
	c    echo.Context
}

func newFriendListServiceTest(t *testing.T) *friendListServiceTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	db, mock := testutil.NewSQLMock(t)
	flr := mock_repository.NewMockFriendListRepository(ctrl)

	return &friendListServiceTest{
		db:   db,
		mock: mock,
		flr:  flr,
		fls:  NewFriendListService(flr),
		c:    testutil.SetUpContextWithDefault(),
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

func Test_friendListService_CheckUserExist(t *testing.T) {
	userId := testutil.UserIDForDebug
	tests := []struct {
		name    string
		expects func(*friendListServiceTest)
		want    bool
		wantErr bool
	}{
		{
			name: "ok: user exist",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(userId).Return(true, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ok: user not exist",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(userId).Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ng: error at CheckUserExist()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(userId).Return(false, testutil.ErrTest)
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			got, err := st.fls.CheckUserExist(st.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListService_GetFriendListByUserId(t *testing.T) {
	userId := testutil.UserIDForDebug
	blockUsers := []int{0}
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(test *friendListServiceTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok: no block user",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(nil, nil)
				st.flr.EXPECT().GetFriendListByUserId(userId).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ok: block some users",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(blockUsers, nil)
				st.flr.EXPECT().GetFriendListByUserIdExcludingBlockUsers(userId, blockUsers).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: error at GetBlockUsersIdList()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(nil, nil)
				st.flr.EXPECT().GetFriendListByUserId(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListByUserIdExcludingBlockUsers()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(blockUsers, nil)
				st.flr.EXPECT().GetFriendListByUserIdExcludingBlockUsers(userId, blockUsers).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			got, err := st.fls.GetFriendListByUserId(st.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListService_GetFriendListOfFriendsByUserId(t *testing.T) {
	userId := testutil.UserIDForDebug
	userList := []int{0}
	userLists := append(userList, userList...)
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(test *friendListServiceTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetFriendListOfFriendsByUserId(userId, userLists).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ok: no 1hop friend",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(nil, nil)
			},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name: "ng: error at GetOneHopFriendsUserIdList()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetBlockUsersIdList()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetFriendListOfFriendsByUserId(userId, userLists).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			got, err := st.fls.GetFriendListOfFriendsByUserId(st.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListService_GetFriendListOfFriendsByUserIdWithPaging(t *testing.T) {
	userId := testutil.UserIDForDebug
	userList := []int{0}
	userLists := append(userList, userList...)
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(test *friendListServiceTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(st *friendListServiceTest) {
				st.c.Set("limit", 0)
				st.c.Set("offset", 0)
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(userId, userLists, 0, 0).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ok: no 1hop friend",
			expects: func(st *friendListServiceTest) {
				st.c.Set("limit", 0)
				st.c.Set("offset", 0)
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(nil, nil)
			},
			want: &model.FriendList{
				Friends: []*model.User(nil),
			},
			wantErr: false,
		},
		{
			name: "ng: error at GetOneHopFriendsUserIdList()",
			expects: func(st *friendListServiceTest) {
				st.c.Set("limit", 0)
				st.c.Set("offset", 0)
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetBlockUsersIdList()",
			expects: func(st *friendListServiceTest) {
				st.c.Set("limit", 0)
				st.c.Set("offset", 0)
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(st *friendListServiceTest) {
				st.c.Set("limit", 0)
				st.c.Set("offset", 0)
				st.flr.EXPECT().GetOneHopFriendsUserIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetBlockUsersIdList(userId).Return(userList, nil)
				st.flr.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(userId, userLists, 0, 0).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			got, err := st.fls.GetFriendListOfFriendsByUserIdWithPaging(st.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserIdWithPaging() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
