package testutil

import (
	"testing"
)

func Test_httpError_SchemaValidation(t *testing.T) {
	tests := []struct {
		schemaName string
		target     any
	}{
		{
			schemaName: "id",
			target:     UserIDForDebug,
		},
		{
			schemaName: "name",
			target:     UserNameForDebug,
		},
	}

	ot := NewOpenAPITester(t, "../../../spec/openapi.yaml")
	for _, tt := range tests {
		t.Run(tt.schemaName, func(t *testing.T) {
			ot.ValidateBySchema(t, tt.schemaName, tt.target)
		})
	}
}
