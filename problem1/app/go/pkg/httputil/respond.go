package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RespondError(c echo.Context, err error) error {
	log.Println(err.Error())

	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return c.JSON(httpErr.StatusCode(), httpErr)
	} else {
		return c.JSON(http.StatusInternalServerError, err)
	}
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
