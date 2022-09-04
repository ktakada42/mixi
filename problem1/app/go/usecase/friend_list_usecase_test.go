package usecase

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"problem1/httputil"
	"problem1/mock/mock_service"
	"problem1/model"
	"problem1/testutil"
)

type friendListUseCaseTest struct {
	db        *sql.DB
	mock      sqlmock.Sqlmock
	fls       *mock_service.MockFriendListService
	flu       FriendListUseCase
	fluStruct *friendListUseCase
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
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(true, nil)
			},
			wantErr: false,
		},
		{
			name: "ng: error at CheckUserExitst()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(false, testutil.ErrTest)
			},
			wantErr: true,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(false, nil)
			},
			wantErr:     true,
			wantErrCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ut := newFriendListUseCaseTest(t)
			tt.expects(ut)

			c, err := httputil.SetUpContext("")
			if err != nil {
				t.Fatal(err)
			}

			err = ut.fluStruct.checkUserExist(c)
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

func Test_friendListUseCase_GetFriendListByUesrId(t *testing.T) {
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
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(true, nil)
				ut.fls.EXPECT().GetFriendListByUserId(gomock.Any()).Return(want, nil)
			},
			want:    want,
			wantErr: false,
		},
		{
			name: "ng: user not exist",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(false, testutil.ErrTest)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(ut *friendListUseCaseTest) {
				ut.fls.EXPECT().CheckUserExist(gomock.Any()).Return(true, nil)
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

			c, err := httputil.SetUpContext("")
			if err != nil {
				t.Fatal(err)
			}

			got, err := ut.flu.GetFriendListByUserId(c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetFriendListByUserId() error = %v, wantErr = %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
