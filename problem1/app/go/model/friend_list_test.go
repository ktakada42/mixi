package model

import (
	"testing"

	"problem1/testutil"
)

func Test_FriendList_SchemaValidation(t *testing.T) {
	target := &FriendList{
		Id:   testutil.UserIDForDebug,
		Name: testutil.UserNameForDebug,
	}

	ot := testutil.NewOpenAPITester(t, "../../../spec/openapi.yaml")
	t.Run(`"FriendList"がOpenAPIの"FriendList"のスキーマと一致している`, func(t *testing.T) {
		ot.ValidateBySchema(t, "FriendList", target)
	})
}
