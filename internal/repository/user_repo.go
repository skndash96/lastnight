package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int32) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	CreateUser(ctx context.Context, name, email string) (db.User, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(d db.DBTX) UserRepository {
	return &userRepository{
		q: db.New(d),
	}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.q.GetUserByEmail(ctx, email)
}

func (r *userRepository) CreateUser(ctx context.Context, name, email string) (db.User, error) {
	return r.q.CreateUser(ctx, db.CreateUserParams{
		Name:  name,
		Email: email,
	})
}
