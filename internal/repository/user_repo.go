package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(d db.DBTX) *UserRepository {
	return &UserRepository{
		q: db.New(d),
	}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email string) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
}
