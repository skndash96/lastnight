package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/helpers"
	"github.com/skndash96/lastnight-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, name, email, password string) (string, error)
}

type authService struct {
	db            *pgxpool.Pool
	tokenProvider auth.TokenProvider
}

func NewAuthService(db *pgxpool.Pool, tokenProvider auth.TokenProvider) AuthService {
	return &authService{
		db:            db,
		tokenProvider: tokenProvider,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	userRepo := repository.NewUserRepository(s.db)
	accountRepo := repository.NewAccountRepository(s.db)

	u, err := userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", NewSrvError(err, SrvErrInvalidInput, "invalid credentials")
		}
		return "", NewSrvError(err, SrvErrInternal, "something went wrong")
	}

	accounts, err := accountRepo.GetUserAccountsByID(ctx, u.ID)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to get user accounts")
	}

	var acc *db.Account
	for _, a := range accounts {
		if a.Provider == "local" {
			acc = &a
			break
		}
	}

	if acc == nil {
		return "", NewSrvError(err, SrvErrInvalidInput, "local account not found")
	}

	if err := bcrypt.CompareHashAndPassword(acc.Password, []byte(password)); err != nil {
		return "", NewSrvError(err, SrvErrInvalidInput, "invalid credentials")
	}

	token, err := s.tokenProvider.GenerateToken(acc.UserID, u.Email)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to generate token")
	}

	return token, nil
}

func (s *authService) Register(ctx context.Context, name, email, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to hash password")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	userRepo := repository.NewUserRepository(tx)
	accountRepo := repository.NewAccountRepository(tx)

	user, err := userRepo.CreateUser(ctx, name, email)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return "", NewSrvError(err, SrvErrInvalidInput, "email already in use")
		}
		return "", NewSrvError(err, SrvErrInternal, "failed to create new user")
	}

	acc, err := accountRepo.CreateAccount(ctx, user.ID, "local", email, string(passwordHash))
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to create new account")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to commit transaction")
	}

	token, err := s.tokenProvider.GenerateToken(acc.UserID, user.Email)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to generate token")
	}

	return token, nil
}
