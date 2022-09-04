package model

import (
	"testing"

	"problem1/testutil"
)

func Test_FriendList_SchemaValidation(t *testing.T) {
	target := &User{
		Id:   testutil.UserIDForDebug,
		Name: testutil.UserNameForDebug,
	}

	ot := testutil.NewOpenAPITester(t, "../../../spec/openapi.yaml")
	t.Run(`"User"がOpenAPIの"User"のスキーマと一致している`, func(t *testing.T) {
		ot.ValidateBySchema(t, "User", target)
	})
}
