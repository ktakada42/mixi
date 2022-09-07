package middleware

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func PagingFunc(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			limit = 20
		}
		if limit < 1 {
			limit = 1
		}
		if limit > 100 {
			limit = 100
		}
		c.Set("limit", limit)

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			page = 1
		}
		if page < 1 {
			page = 1
		}
		c.Set("page", page)

		offset := limit * (page - 1)
		c.Set("offset", offset)

		return next(c)
	}
}
