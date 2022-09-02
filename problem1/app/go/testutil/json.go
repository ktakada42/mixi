package testutil

import (
	"bytes"
	"encoding/json"
	"testing"
)

func I2V(t *testing.T, src, dst interface{}) {
	t.Helper()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(&buf).Decode(&dst); err != nil {
		t.Fatal(err)
	}
}
