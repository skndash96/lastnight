package app

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/handler"
	"github.com/skndash96/lastnight-backend/internal/provider"
	"github.com/skndash96/lastnight-backend/internal/queue"
	"github.com/skndash96/lastnight-backend/internal/repository"
	"github.com/skndash96/lastnight-backend/internal/service"
)

func RegisterRoutes(e *echo.Echo, cfg *config.AppConfig, pool *pgxpool.Pool, ingestionQ *queue.IngestionQ) {
	r := e.Group("/api")

	authRepo := repository.NewAuthRepository(pool)
	teamRepo := repository.NewTeamRepository(pool)

	sessionProvider := provider.NewSessionProvider(cfg.Auth.Session, authRepo)
	uploadProvider, err := provider.NewUploadProvider(cfg.Minio)

	authSrv := service.NewAuthService(pool, sessionProvider)
	docSrv := service.NewDocumentService(uploadProvider, pool, ingestionQ)
	teamSrv := service.NewTeamService(pool)
	tagSrv := service.NewTagService(pool)

	if err != nil {
		log.Fatalf("failed to initialize upload provider: %v", err)
	}

	r.Use(auth.SessionMW(sessionProvider, cfg.Auth.Cookie))

	{
		h := handler.NewHealthHandler()
		g := r.Group("/health")
		g.GET("", h.HealthCheck)
	}

	{
		h := handler.NewAuthHandler(cfg, authSrv)
		g := r.Group("/auth")
		g.POST("/login", h.Login)
		g.POST("/register", h.Register)
		g.DELETE("/logout", h.Logout)
	}

	{
		team_h := handler.NewTeamHandler(teamSrv)
		tag_h := handler.NewTagHandler(tagSrv)

		teamsG := r.Group("/teams")

		teamsG.GET("", team_h.GetTeams)
		teamsG.POST("/default", team_h.JoinDefaultTeam)

		teamG := teamsG.Group("/:teamID")
		teamG.Use(auth.TeamMW(teamRepo))

		teamG.GET("/filters", tag_h.ListFilters)
		teamG.PUT("/filters", tag_h.UpdateFilters)

		{
			tagsG := teamG.Group("/tags")

			tagsG.POST("", tag_h.CreateTagKey)
			tagsG.PATCH("/:tagID", tag_h.UpdateTagKey)
			tagsG.DELETE("/:tagID", tag_h.DeleteTagKey)

			tagsG.POST("/:tagID/values", tag_h.CreateTagValue)
			tagsG.DELETE("/:tagID/values/:tagValueID", tag_h.DeleteTagValue)
		}

		{
			h := handler.NewUploadHandler(docSrv)
			g := teamG.Group("/uploads")
			g.POST("/presign", h.PresignUpload)
			g.POST("/commit", h.CommitUpload)
		}

		{
			h := handler.NewDocHandler(docSrv)

			docRefsG := teamG.Group("/refs")

			docRefsG.PUT("/:docRefID/tags", h.UpdateDocRefTags)
		}
	}

	// for workers
	{
		g := r.Group("/internal")
		g.Use(auth.InternalMw(cfg.WorkerToken))

		internalH := handler.NewInternalHandler(docSrv)

		g.GET("/docs/:docID", internalH.GetDoc)
		g.PATCH("/docs/:docID/status", internalH.UpdateDocProcStatus)
	}
}
