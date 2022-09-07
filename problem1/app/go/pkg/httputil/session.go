package httputil

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

func SetUpContext(rawURL string) (echo.Context, error) {
	e := echo.New()
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return e.NewContext(&http.Request{URL: url}, nil), nil
}
