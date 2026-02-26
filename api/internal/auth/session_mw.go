package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/config"
)

const authKeyInCtx = "lastnight_actor"

func SessionMW(tokenProvider TokenProvider, cookieCfg config.CookieConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Request().Cookie(cookieCfg.Name)
			if err == nil && token.Value != "" {
				actor, err := tokenProvider.ValidateToken(c.Request().Context(), token.Value)
				if err == nil {
					c.Set(authKeyInCtx, actor)
				}
			}

			return next(c)
		}
	}
}

func GetSession(c echo.Context) (s *Actor, ok bool) {
	s, ok = c.Get(authKeyInCtx).(*Actor)

	if !ok || s.UserID == 0 {
		return nil, false
	}

	return s, true
}
