package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"problem1/pkg/testutil"
)

func Test_respond_RespondError(t *testing.T) {
	tests := []struct {
		name     string
		wantCode int
		err      error
	}{
		{
			name:     "HTTPError",
			wantCode: http.StatusInternalServerError,
			err:      NewHTTPError(testutil.ErrTest, http.StatusInternalServerError, ""),
		},
		{
			name:     "normalError",
			wantCode: http.StatusInternalServerError,
			err:      testutil.ErrTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(nil, rec)
			RespondError(c, tt.err)
			resp := rec.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

type friend struct {
	userId int    `json:"userId"`
	name   string `json:"name"`
}

type friendList struct {
	friends []*friend `json:"friends"`
}

func Test_respond_RespondJSON(t *testing.T) {
	tests := []struct {
		name     string
		payload  any
		want     string
		wantCode int
		wantErr  bool
	}{
		{
			name: "ok",
			payload: &friendList{
				friends: []*friend{
					{
						userId: testutil.UserIDForDebug,
						name:   testutil.UserNameForDebug,
					},
					{
						userId: testutil.UserIDForDebug,
						name:   testutil.UserNameForDebug,
					},
				},
			},
			wantCode: http.StatusOK,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(nil, rec)
			RespondJSON(c, tt.wantCode, tt.payload)
			resp := rec.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.wantCode, resp.StatusCode)
			if tt.wantErr {
				// エラー時はBodyを書き換えているのでチェック
				testutil.AssertResponseBody(t, tt.want, resp.Body)
			}
		})
	}
}
