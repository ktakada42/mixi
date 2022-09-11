package httputil

import (
	"errors"
	"fmt"
)

type HTTPError interface {
	error
	StatusCode() int
}

type httpError struct {
	origin     error
	statusCode int
	message    string
}

func NewHTTPError(origin error, statusCode int, message string) error {
	if message == "" {
		message = origin.Error()
	}

	return &httpError{
		origin:     origin,
		statusCode: statusCode,
		message:    message,
	}
}

func (e *httpError) Error() string {
	return fmt.Sprintf("StatusCode = %d, msg = %s", e.statusCode, e.message)
}

func (e *httpError) StatusCode() int {
	return e.statusCode
}

func As(err error, c int) bool {
	var hErr HTTPError
	if errors.As(err, &hErr) && hErr.StatusCode() == c {
		return true
	}

	return false
}
