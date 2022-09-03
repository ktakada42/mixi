package httputil

import (
	"io"
	"net/http"
	"net/http/httptest"
)

func NewRequestAndRecorder(method, path string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest(method, path, body)
}
