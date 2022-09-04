package usecase

import (
	"database/sql"
	"net/http"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/mock/mock_service"
	"problem1/model"
	"problem1/testutil"
)

type friendListUseCaseTest struct {
	db   *sql.DB
	mock sqlmock.Sqlmock
	fls  *mock_service.MockFriendListService
	flu  FriendListUseCase
}

func newFriendListUseCaseTest(t *testing.T) *friendListUseCaseTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	db, mock := testutil.NewSQLMock(t)
	fls := mock_service.NewMockFriendListService(ctrl)
	flu := NewFriendListUseCase(db, fls)

	return &friendListUseCaseTest{
		db:   db,
		mock: mock,
		fls:  fls,
		flu:  flu,
	}
}

func newFriendList() []*model.User {
	return []*model.User{
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

func Test_friendListUseCase_GetFriendListByUesrId(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(test *friendListUseCaseTest)
		want    []*model.User
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().GetFriendListByUserId(gomock.Any()).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().GetFriendListByUserId(gomock.Any()).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			e := echo.New()
			e.GET("", nil)
			url, err := url.Parse("")
			if err != nil {
				t.Fatal(err)
			}
			c := e.NewContext(&http.Request{URL: url}, nil)

			got, err := ut.flu.GetFriendListByUserId(c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
