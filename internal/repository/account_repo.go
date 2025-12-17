package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type AccountRepository interface {
	GetUserAccountsByID(ctx context.Context, userID int32) ([]db.Account, error)
	CreateAccount(ctx context.Context, userID int32, provider, providerAccountID, password string) (db.Account, error)
}

type accountRepository struct {
	q *db.Queries
}

func NewAccountRepository(d db.DBTX) AccountRepository {
	return &accountRepository{
		q: db.New(d),
	}
}

func (r *accountRepository) GetUserAccountsByID(ctx context.Context, userID int32) ([]db.Account, error) {
	return r.q.GetUserAccountsByID(ctx, userID)
}

func (r *accountRepository) CreateAccount(ctx context.Context, userID int32, provider, providerAccountID, password string) (db.Account, error) {
	return r.q.CreateAccount(ctx, db.CreateAccountParams{
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: providerAccountID,
		Password:          []byte(password),
	})
}
