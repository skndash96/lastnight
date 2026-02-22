package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type authHandler struct {
	appCfg *config.AppConfig
	srv    service.AuthService
}

func NewAuthHandler(appCfg *config.AppConfig, srv service.AuthService) *authHandler {
	return &authHandler{
		appCfg: appCfg,
		srv:    srv,
	}
}

// @Summary Login endpoint
// @Tags Auth
// @Description Authenticates a user with email and password
// @Param user body dto.LoginRequest true "Login credentials"
// @Produce json
// @Success 200
// @Failure default	{object} dto.ErrorResponse
// @Router /api/auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var body dto.LoginRequest
	if err := c.Bind(&body); err != nil {
		return err
	}

	if err := c.Validate(&body); err != nil {
		a, ok := err.(*validator.ValidationErrors)
		fmt.Printf("Validation error: %+v, ok: %v", a, ok)
		return err
	}

	token, err := h.srv.Login(c.Request().Context(), body.Email, body.Password)
	if err != nil {
		return err
	}

	c.SetCookie(auth.NewCookie(h.appCfg.Auth.Cookie, token))

	c.NoContent(http.StatusOK)

	return nil
}

// @Summary Signup endpoint
// @Tags Auth
// @Description Registers a new user and returns an authentication token
// @Param user body dto.RegisterRequest true "Register credentials"
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/auth/register [post]
func (h *authHandler) Register(c echo.Context) error {
	var body dto.RegisterRequest
	if err := c.Bind(&body); err != nil {
		return err
	}

	if err := c.Validate(&body); err != nil {
		return err
	}

	token, err := h.srv.Register(c.Request().Context(), body.Name, body.Email, body.Password)
	if err != nil {
		return err
	}

	c.SetCookie(auth.NewCookie(h.appCfg.Auth.Cookie, token))

	c.NoContent(http.StatusOK)

	return nil
}

// @Summary Logout endpoint
// @Tags Auth
// @Description Logs out the current user by clearing the authentication token cookie
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/auth/logout [delete]
func (h *authHandler) Logout(c echo.Context) error {
	if _, ok := auth.GetSession(c); !ok {
		return echo.ErrUnauthorized
	}

	c.SetCookie(&http.Cookie{
		Name:     h.appCfg.Auth.Cookie.Name,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	c.NoContent(http.StatusOK)

	return nil
}
