package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/skndash96/lastnight-backend/internal/config"
)

type jwtProvider struct {
	cfg config.JWTConfig
}

func NewJwtProvider(jwtCfg config.JWTConfig) TokenProvider {
	return &jwtProvider{
		cfg: jwtCfg,
	}
}

func (s *jwtProvider) GenerateToken(userID int32, email string) (string, error) {
	claims := AuthClaims{
		AuthIdentity: AuthIdentity{
			UserID: userID,
			Email:  email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.Expiry)),
			Issuer:    s.cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.cfg.Secret)
}

func (s *jwtProvider) ParseToken(tokenStr string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenUnverifiable
		}
		return s.cfg.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
