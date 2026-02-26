package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type teamHandler struct {
	teamSrv *service.TeamService
}

func NewTeamHandler(teamSrv *service.TeamService) *teamHandler {
	return &teamHandler{
		teamSrv: teamSrv,
	}
}

// @Summary Get Teams
// @Tags Team
// @Description Get the user's teams list
// @Produce json
// @Success 200 {object} dto.GetTeamsResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams [get]
func (h *teamHandler) GetTeams(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	ts, err := h.teamSrv.GetTeamsByUserID(c.Request().Context(), session.UserID)
	if err != nil {
		return err
	}

	c.JSON(200, &dto.GetTeamsResponse{
		Data: ts,
	})

	return nil
}

// @Summary Join Default Team
// @Tags Team
// @Description Join the default team for the user
// @Produce json
// @Success 201 {object} dto.JoinTeamResponse
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

	c.JSON(http.StatusCreated, &dto.JoinTeamResponse{
		Data: p,
	})

	return nil
}
