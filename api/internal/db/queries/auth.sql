-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;



-- name: GetUserAccountsByID :many
SELECT * FROM accounts WHERE user_id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (user_id, provider, provider_account_id, password) VALUES ($1, $2, $3, $4) RETURNING *;



-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1;

-- name: CreateSession :one
INSERT INTO sessions (user_id, email, expiry) VALUES ($1, $2, $3) RETURNING *;
