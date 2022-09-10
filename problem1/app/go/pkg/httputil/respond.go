package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"

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
	v := reflect.ValueOf(payload)
	if v.Kind() == reflect.Slice {
		// any -> []any に変換
		payloads := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			payloads[i] = v.Index(i).Interface()
		}

		respondJSON(c, status, payloads)
		return
	}

	respondJSON(c, status, payload)
}

func respondJSON(c echo.Context, status int, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		RespondError(c, err)
	}

	w := c.Response().Writer
	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		RespondError(c, err)
	}
}
