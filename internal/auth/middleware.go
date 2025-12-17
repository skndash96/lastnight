package auth

import (
	"github.com/labstack/echo/v4"
)

const authKeyInCtx = "lastnight_token"

func AuthMW(tokenProvider TokenProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Request().Cookie(authKeyInCtx)
			if err == nil && token.Value != "" {
				claims, err := tokenProvider.ParseToken(token.Value)
				if err == nil {
					c.Set(authKeyInCtx, claims.AuthIdentity)
				}
			}

			return next(c)
		}
	}
}

func GetSession(c echo.Context) (s AuthIdentity, ok bool) {
	s, ok = c.Get(authKeyInCtx).(AuthIdentity)

	if !ok || s.Email == "" {
		return AuthIdentity{}, false
	}

	return s, true
}
