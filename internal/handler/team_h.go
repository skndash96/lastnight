package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type TeamHandler interface {
	GetDefaultTeam(c echo.Context) error
	JoinDefaultTeam(c echo.Context) error
}

type teamHandler struct {
	teamSrv *service.TeamService
}

func NewTeamHandler(teamSrv *service.TeamService) TeamHandler {
	return &teamHandler{
		teamSrv: teamSrv,
	}
}

// @Summary Get Default Team
// @Tags Team
// @Description Get the default team for the user
// @Produce json
// @Success default {object} dto.GetTeamResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/default [get]
func (h *teamHandler) GetDefaultTeam(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	p, err := h.teamSrv.GetDefaultTeam(c.Request().Context(), session.Email)
	if err != nil {
		return err
	}

	c.JSON(200, &dto.GetTeamResponse{
		Data: p,
	})

	return nil
}

// @Summary Join Default Team
// @Tags Team
// @Description Join the default team for the user
// @Produce json
// @Success default {object} dto.JoinTeamResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/default [post]
func (h *teamHandler) JoinDefaultTeam(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	p, err := h.teamSrv.JoinDefaultTeam(c.Request().Context(), session.UserID, session.Email)
	if err != nil {
		return err
	}

	c.JSON(200, &dto.JoinTeamResponse{
		Data: p,
	})

	return nil
}
