package usecase

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/mock/mock_service"
	"problem1/model"
	"problem1/pkg/httputil"
	"problem1/pkg/testutil"
)

type friendListUseCaseTest struct {
	db        *sql.DB
	mock      sqlmock.Sqlmock
	fls       *mock_service.MockFriendListService
	flu       FriendListUseCase
	fluStruct *friendListUseCase
	c         echo.Context
}

func newFriendListUseCaseTest(t *testing.T) *friendListUseCaseTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	db, mock := testutil.NewSQLMock(t)
	fls := mock_service.NewMockFriendListService(ctrl)
	flu := NewFriendListUseCase(db, fls)

	return &friendListUseCaseTest{
		db:        db,
		mock:      mock,
		fls:       fls,
		flu:       flu,
		fluStruct: flu.(*friendListUseCase),
		c:         testutil.SetUpContextWithDefault(),
	}
}

func newFriendList() *model.FriendList {
	return &model.FriendList{
		Friends: []*model.Friend{
			{
				UserId: 111111,
				Name:   "hoge",
			},
			{
				UserId: 222222,
				Name:   "fuga",
			},
		},
	}
}

func Test_friendListUseCase_checkUserExist(t *testing.T) {
	tests := []struct {
		name        string
		expects     func(*friendListUseCaseTest)
		wantErr     bool
		wantErrCode int
	}{
		{
			name: "ok",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name: "ng: error at CheckUserExist()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(false, testutil.ErrTest)
			},
			wantErr: true,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(false, nil)
			},
			wantErr:     true,
			wantErrCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			err := ut.fluStruct.checkUserExist(testutil.UserIDForDebug)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkUserExist() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if err != nil && tt.wantErrCode != 0 {
				if !httputil.As(err, tt.wantErrCode) {
					t.Fatalf("checkUserExist() error = %v, wantErrCode= %v", err, tt.wantErrCode)
				}
			}
		})
	}
}

func Test_friendListUseCase_GetFriendListByUserId(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(*friendListUseCaseTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListByUserId(ut.c).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(false, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListByUserId(ut.c).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			got, err := ut.flu.GetFriendListByUserId(ut.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListUseCase_GetFriendListOfFriendsByUserId(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(*friendListUseCaseTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListOfFriendsByUserId(ut.c).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(false, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListOfFriendsByUserId(ut.c).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			got, err := ut.flu.GetFriendListOfFriendsByUserId(ut.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_friendListUseCase_GetFriendListOfFriendsByUserIdWithPaging(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name    string
		expects func(*friendListUseCaseTest)
		want    *model.FriendList
		wantErr bool
	}{
		{
			name: "ok",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(ut.c).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(false, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(testutil.UserIDForDebug).Return(true, nil)
				ut.fls.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(ut.c).Return(nil, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			got, err := ut.flu.GetFriendListOfFriendsByUserIdWithPaging(ut.c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListOfFriendsByUserIdWithPaging() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
