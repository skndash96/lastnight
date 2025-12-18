package config

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	IsProd bool
	Port   int
	DbURL  string
	Auth   AuthConfig
}

type AuthConfig struct {
	Session SessionConfig
	Cookie  CookieConfig
}

type SessionConfig struct {
	Expiry time.Duration
}

type CookieConfig struct {
	Name     string
	Secure   bool
	SameSite http.SameSite
}

func New() *AppConfig {
	port, err := strconv.Atoi(GetEnv("PORT", "1323"))
	if err != nil {
		port = 1323
	}

	isProd, err := strconv.ParseBool(GetEnv("IS_PROD", "true"))
	if err != nil {
		isProd = false
	}

	appCfg := &AppConfig{
		IsProd: isProd,
		Port:   port,
		DbURL:  GetEnv("GOOSE_DBSTRING", ""),

		Auth: AuthConfig{
			Session: SessionConfig{
				Expiry: time.Duration(14*24) * time.Hour,
			},
			Cookie: CookieConfig{
				Name:     "lastnight_token",
				Secure:   isProd,
				SameSite: http.SameSiteLaxMode,
			},
		},
	}

	return appCfg
}

func GetEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
