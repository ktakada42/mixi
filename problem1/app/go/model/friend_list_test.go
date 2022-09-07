package model

import (
	"testing"

	"problem1/pkg/testutil"
)

func Test_FriendList_SchemaValidation(t *testing.T) {
	tests := []struct {
		schemaName string
		target     any
	}{
		{
			schemaName: "User",
			target: &User{
				Id:   testutil.UserIDForDebug,
				Name: testutil.UserNameForDebug,
			},
		},
		{
			schemaName: "FriendList",
			target: &FriendList{
				Friends: []*User{
					{
						Id:   testutil.UserIDForDebug,
						Name: testutil.UserNameForDebug,
					},
				},
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
