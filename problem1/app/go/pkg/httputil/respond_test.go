package httputil

import (
	"net/http"
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
			rec, req := NewRequestAndRecorder("GET", "/test", testutil.I2Reader(t, "body"))
			e := echo.New()
			e.GET("/test", func(c echo.Context) error {
				return RespondError(c, tt.err)
			})
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}
