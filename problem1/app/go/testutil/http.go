package testutil

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertResponseBody(t *testing.T, want any, body io.Reader) {
	t.Helper()

	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := io.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, string(b), string(got))
}
