package httputil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"problem1/pkg/testutil"
)

func Test_httpError_NewHTTPError(t *testing.T) {
	statusCode := http.StatusInternalServerError

	tests := []struct {
		name   string
		origin error
		want   error
	}{
		{
			name:   "ok",
			origin: testutil.ErrTest,
			want: &httpError{
				origin:     testutil.ErrTest,
				statusCode: statusCode,
				message:    testutil.ErrTest.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewHTTPError(tt.origin, statusCode, "")
			assert.Equal(t, tt.want, err)
		})
	}
}

func Test_httpError_As(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ok",
			err:  NewHTTPError(testutil.ErrTest, http.StatusInternalServerError, ""),
			want: true,
		},
		{
			name: "ng: status code not equal",
			err:  NewHTTPError(testutil.ErrTest, http.StatusBadRequest, ""),
			want: false,
		},
		{
			name: "ng: not httpError",
			err:  testutil.ErrTest,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := As(tt.err, http.StatusInternalServerError)
			assert.Equal(t, tt.want, got)
		})
	}
}
