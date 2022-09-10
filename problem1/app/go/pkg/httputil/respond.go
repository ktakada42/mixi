package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RespondError(c echo.Context, err error) {
	w := c.Response().Writer

	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		w.WriteHeader(httpErr.StatusCode())
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println(err.Error())
}

func RespondJSON(c echo.Context, status int, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		RespondError(c, err)
		return
	}

	w := c.Response().Writer
	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		RespondError(c, err)
	}
}
