package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/skndash96/lastnight-backend/internal/db"
)

type AuthRepository struct {
	q *db.Queries
}

func NewAuthRepository(d db.DBTX) *AuthRepository {
	return &AuthRepository{
		q: db.New(d),
	}
}

// -------- user --------
func (r *AuthRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *AuthRepository) CreateUser(ctx context.Context, name, email string) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
}

// -------- account --------
func (r *AuthRepository) GetUserAccountsByID(ctx context.Context, userID int32) ([]db.Account, error) {
	return r.q.GetUserAccountsByID(ctx, userID)
}

func (r *AuthRepository) CreateAccount(ctx context.Context, userID int32, provider, providerAccountID, password string) (db.Account, error) {
	return r.q.CreateAccount(ctx, db.CreateAccountParams{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
		Password:          []byte(password),
	})
}

// -------- session --------
func (r *AuthRepository) GetSessionByID(ctx context.Context, id pgtype.UUID) (db.Session, error) {
	return r.q.GetSessionByID(ctx, id)
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID int32, email string, expiry time.Time) (*db.Session, error) {
	session, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		UserID: userID,
		Email:  email,
		Expiry: expiry,
	})
	if err != nil {
		return nil, err
	}

	return &session, nil
}
