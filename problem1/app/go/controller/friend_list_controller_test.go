package controller

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/mock/mock_usecase"
	"problem1/model"
	"problem1/pkg/httputil"
	"problem1/pkg/testutil"
)

type friendListControllerTest struct {
	flu  *mock_usecase.MockFriendListUseCase
	flc  FriendListController
	echo *echo.Echo
}

func newFriendListControllerTest(t *testing.T) *friendListControllerTest {
	t.Helper()

	ctrl := gomock.NewController(t)
	flu := mock_usecase.NewMockFriendListUseCase(ctrl)

	return &friendListControllerTest{
		flu:  flu,
		flc:  NewFriendListController(flu),
		echo: echo.New(),
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

func Test_friendListController_GetFriendListByUserId(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name       string
		expects    func(test *friendListControllerTest)
		url        string
		want       *model.FriendList
		wantStatus int
		wantErr    bool
	}{
		{
			name: "ok",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListByUserId(gomock.Any()).Return(want, nil)
			},
			url:        "/get_friend_list?ID=123456789",
			want:       want,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "ng: userId missing in query parameter",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId not integer",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=invalid",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId minus",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=-1",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId over mysql max limit",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=999999999999999999",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "ng: error at GetFriendListByUserId()",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListByUserId(gomock.Any()).Return(nil, testutil.ErrTest)
			},
			url:        "/get_friend_list?ID=0",
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := newFriendListControllerTest(t)
			tt.expects(ct)

			rec, req := httputil.NewRequestAndRecorder("GET", tt.url, nil)
			ct.echo.GET("/get_friend_list", ct.flc.GetFriendListByUserId)
			ct.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if !tt.wantErr {
				testutil.AssertResponseBody(t, want, rec.Body)
			}
		})
	}
}

func Test_friendListController_GetFriendListOfFriendsByUserId(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name       string
		expects    func(test *friendListControllerTest)
		url        string
		want       *model.FriendList
		wantStatus int
		wantErr    bool
	}{
		{
			name: "ok",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListOfFriendsByUserId(gomock.Any()).Return(want, nil)
			},
			url:        "/get_friend_list?ID=123456789",
			want:       want,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "ng: userId missing in query parameter",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId not integer",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=invalid",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId minus",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=-1",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId over mysql max limit",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=999999999999999999",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListOfFriendsByUserId(gomock.Any()).Return(nil, testutil.ErrTest)
			},
			url:        "/get_friend_list?ID=0",
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := newFriendListControllerTest(t)
			tt.expects(ct)

			rec, req := httputil.NewRequestAndRecorder("GET", tt.url, nil)
			ct.echo.GET("/get_friend_list", ct.flc.GetFriendListOfFriendsByUserId)
			ct.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if !tt.wantErr {
				testutil.AssertResponseBody(t, want, rec.Body)
			}
		})
	}
}

func Test_friendListController_GetFriendListOfFriendsByUserIdWithPaging(t *testing.T) {
	want := newFriendList()

	tests := []struct {
		name       string
		expects    func(test *friendListControllerTest)
		url        string
		want       *model.FriendList
		wantStatus int
		wantErr    bool
	}{
		{
			name: "ok",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(gomock.Any()).Return(want, nil)
			},
			url:        "/get_friend_list?ID=123456789",
			want:       want,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "ng: userId missing in query parameter",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId not integer",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=invalid",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId minus",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=-1",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "ng: userId over mysql max limit",
			expects:    func(ct *friendListControllerTest) {},
			url:        "/get_friend_list?ID=999999999999999999",
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "ng: error at GetFriendListOfFriendsByUserId()",
			expects: func(ct *friendListControllerTest) {
				ct.flu.EXPECT().GetFriendListOfFriendsByUserIdWithPaging(gomock.Any()).Return(nil, testutil.ErrTest)
			},
			url:        "/get_friend_list?ID=0",
			want:       nil,
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := newFriendListControllerTest(t)
			tt.expects(ct)

			rec, req := httputil.NewRequestAndRecorder("GET", tt.url, nil)
			ct.echo.GET("/get_friend_list", ct.flc.GetFriendListOfFriendsByUserIdWithPaging)
			ct.echo.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			if !tt.wantErr {
				testutil.AssertResponseBody(t, want, rec.Body)
			}
		})
	}
}
