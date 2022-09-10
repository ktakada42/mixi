package httputil

import (
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
