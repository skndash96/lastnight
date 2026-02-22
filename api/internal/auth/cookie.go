package auth

import (
	"net/http"

	"github.com/skndash96/lastnight-backend/internal/config"
)

func NewCookie(c config.CookieConfig, token string) *http.Cookie {
	return &http.Cookie{
		Name:     c.Name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: c.SameSite,
		// no need for max-age, because jwt has expiry
	}
}
