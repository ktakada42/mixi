package httputil

import (
	"net/http"
	"testing"

	"problem1/testutil"
)

func Test_httpError_SchemaValidation(t *testing.T) {
	tests := []struct {
		schemaName string
		target     any
	}{
		{
			schemaName: "HTTPError",
			target: &httpError{
				origin:     testutil.ErrTest,
				statusCode: http.StatusInternalServerError,
				message:    testutil.ErrTest.Error(),
			},
		},
	}

	ot := testutil.NewOpenAPITester(t, "../../../spec/openapi.yaml")
	for _, tt := range tests {
		t.Run(tt.schemaName, func(t *testing.T) {
			ot.ValidateBySchema(t, tt.schemaName, tt.target)
		})
	}
}
