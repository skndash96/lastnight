package auth

import "github.com/golang-jwt/jwt/v5"

type AuthIdentity struct {
	UserID int32  `json:"uid"`
	Email  string `json:"email"`
}

type AuthClaims struct {
	AuthIdentity
	jwt.RegisteredClaims
}

type TokenProvider interface {
	GenerateToken(userID int32, email string) (string, error)
	ParseToken(tokenStr string) (*AuthClaims, error)
}
