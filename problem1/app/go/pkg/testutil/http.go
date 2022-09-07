package testutil

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/labstack/echo/v4"
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

func SetUpContextWithDefault() echo.Context {
	c := echo.New().NewContext(nil, nil)
	c.Set("userId", 123456789)

	return c
}
