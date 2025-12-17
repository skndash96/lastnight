package api

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/handler"
	"github.com/skndash96/lastnight-backend/internal/service"
)

func RegisterRoutes(e *echo.Echo, cfg *config.AppConfig, pool *pgxpool.Pool) {
	r := e.Group("/api")

	jwtProvider := auth.NewJwtProvider(cfg.Auth.JWT)

	r.Use(auth.AuthMW(jwtProvider))

	{
		h := handler.NewHealthHandler()
		g := r.Group("/health")
		g.GET("", h.HealthCheck)
	}

	{
		authSrv := service.NewAuthService(pool, jwtProvider)

		h := handler.NewAuthHandler(cfg, authSrv)
		g := r.Group("/auth")
		g.POST("/login", h.Login)
		g.POST("/register", h.Register)
		g.DELETE("/logout", h.Logout)
	}

	{
		teamSrv := service.NewTeamService(pool)
		h := handler.NewTeamHandler(teamSrv)

		g := r.Group("/teams")

		g.GET("/default", h.GetDefaultTeam)
		g.POST("/default", h.JoinDefaultTeam)
	}
}
