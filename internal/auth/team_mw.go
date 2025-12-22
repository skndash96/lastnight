package auth

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

func TeamMW(teamRepo *repository.TeamRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			identity, ok := GetSession(c)
			if !ok {
				return echo.ErrUnauthorized
			}

			teamID, err := strconv.Atoi(c.Param("teamID"))
			if err != nil || teamID == 0 {
				return echo.ErrBadRequest
			}

			tm, err := teamRepo.GetTeamMembershipByUserID(c.Request().Context(), identity.UserID, int32(teamID))
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return echo.ErrForbidden
				}
				return echo.ErrInternalServerError
			}

			identity.MembershipID = tm.ID
			identity.TeamID = tm.TeamID
			identity.Role = tm.Role

			c.Set(authKeyInCtx, identity)

			return next(c)
		}
	}
}
