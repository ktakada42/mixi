package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

func I2Reader(t *testing.T, i any) io.Reader {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(i); err != nil {
		t.Fatal(err)
	}

	return buf
}
