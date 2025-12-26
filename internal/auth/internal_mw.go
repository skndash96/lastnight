package auth

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func InternalMw(workerToken string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}

			if token != fmt.Sprintf("Bearer %s", workerToken) {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			return next(c)
		}
	}
}
