package testutil

import (
	"bytes"
	"encoding/json"
	"io"
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

func I2Reader(t *testing.T, i any) io.Reader {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(i); err != nil {
		t.Fatal(err)
	}

	return buf
}
