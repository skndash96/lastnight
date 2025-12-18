package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type AccountRepository struct {
	q *db.Queries
}

func NewAccountRepository(d db.DBTX) *AccountRepository {
	return &AccountRepository{
		q: db.New(d),
	}
}

func (r *AccountRepository) GetUserAccountsByID(ctx context.Context, userID int32) ([]db.Account, error) {
	return r.q.GetUserAccountsByID(ctx, userID)
}

func (r *AccountRepository) CreateAccount(ctx context.Context, userID int32, provider, providerAccountID, password string) (db.Account, error) {
	return r.q.CreateAccount(ctx, db.CreateAccountParams{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
		Password:          []byte(password),
	})
}
