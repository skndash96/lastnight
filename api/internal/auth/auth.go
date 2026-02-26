package auth

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type Actor struct {
	Email        string
	UserID       int32
	TeamID       int32
	MembershipID int32
	Role         db.TeamUserRole
}

type TokenProvider interface {
	GenerateToken(ctx context.Context, userID int32, email string) (string, error)
	ValidateToken(ctx context.Context, token string) (*Actor, error)
}
