package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/dto"
)

type healthHandler struct{}

func NewHealthHandler() *healthHandler {
	return &healthHandler{}
}

// @Summary Healthcheck endpoint
// @Tags Health
// @Description Returns the health status of the application
// @Produce json
// @Success 200 {object} dto.HealthResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/health [get]
func (h *healthHandler) HealthCheck(c echo.Context) error {
	c.JSON(http.StatusOK, &dto.HealthResponse{
		Status: "ok",
	})

	return nil
}
