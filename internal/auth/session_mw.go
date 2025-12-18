package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/db"
)

const authKeyInCtx = "lastnight_actor"

type Actor struct {
	AuthIdentity
	TeamID int32
	Role   db.TeamUserRole
}

func SessionMW(tokenProvider TokenProvider, cookieCfg config.CookieConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := c.Request().Cookie(cookieCfg.Name)
			if err == nil && token.Value != "" {
				claims, err := tokenProvider.ParseToken(token.Value)
				if err == nil {
					c.Set(authKeyInCtx, Actor{
						AuthIdentity: claims.AuthIdentity,
						// TeamID and Role, filled by TeamMW
					})
				}
			}

			return next(c)
		}
	}
}

func GetSession(c echo.Context) (s Actor, ok bool) {
	s, ok = c.Get(authKeyInCtx).(Actor)

	if !ok || s.UserID == 0 {
		return Actor{}, false
	}

	return s, true
}
