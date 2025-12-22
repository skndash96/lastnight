package service

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type SrvError struct {
	internal error
	Kind     SrvErrKind
	Message  string
}

type SrvErrKind string

func (e *SrvError) Error() string {
	return e.Message
}

func (e *SrvError) Unwrap() error {
	return e.internal
}

func NewSrvError(err error, kind SrvErrKind, message string) *SrvError {
	var pgErr *pgconn.PgError

	if ok := errors.As(err, &pgErr); ok {
		return mapPgErrToSrvError(pgErr)
	}

	return &SrvError{
		internal: err,
		Kind:     kind,
		Message:  message,
	}
}

func mapPgErrToSrvError(err *pgconn.PgError) *SrvError {
	switch err.Code {
	case "23505", "23503":
		return &SrvError{
			internal: err,
			Kind:     SrvErrConflict,
			Message:  "conflicting values",
		}
	default:
		return &SrvError{
			internal: err,
			Kind:     SrvErrInternal,
			Message:  "unknown error",
		}
	}
}

const (
	// Client fault (input or business logic)
	SrvErrInvalidInput SrvErrKind = "invalid_input" // malformed data, violated constraints
	SrvErrUnauthorized SrvErrKind = "unauthorized"  // auth failed
	SrvErrForbidden    SrvErrKind = "forbidden"     // auth ok, but access denied
	SrvErrNotFound     SrvErrKind = "not_found"     // lookup failed
	SrvErrConflict     SrvErrKind = "conflict"      // duplicate resource / violation (email exists)

	// Non-client, unexpected runtime/system errors
	SrvErrInternal SrvErrKind = "internal_error" // unexpected application failure
)
