package service

import (
	"database/sql"
	"net/http"
	"net/url"
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
	tests := []struct {
		name    string
		expects func(*friendListServiceTest)
		want    bool
		wantErr bool
	}{
		{
			name: "ok: true",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(gomock.Any()).Return(true, nil)
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ok: false",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(gomock.Any()).Return(false, nil)
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ng: error at CheckUserExist()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().CheckUserExist(gomock.Any()).Return(false, testutil.ErrTest)
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			e := echo.New()
			e.GET("", nil)
			url, err := url.Parse("")
			if err != nil {
				t.Fatal(err)
			}
			c := e.NewContext(&http.Request{URL: url}, nil)

			got, err := st.fls.CheckUserExist(c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListService_GetFriendListByUserId(t *testing.T) {
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
				st.flr.EXPECT().GetFriendListByUserId(gomock.Any()).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(st *friendListServiceTest) {
				st.flr.EXPECT().GetFriendListByUserId(gomock.Any()).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := newFriendListServiceTest(t)
			tt.expects(st)

			e := echo.New()
			e.GET("", nil)
			url, err := url.Parse("")
			if err != nil {
				t.Fatal(err)
			}
			c := e.NewContext(&http.Request{URL: url}, nil)

			got, err := st.fls.GetFriendListByUserId(c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
