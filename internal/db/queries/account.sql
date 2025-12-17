-- name: GetUserAccountsByID :many
SELECT * FROM accounts WHERE user_id = $1;

-- name: CreateAccount :one
INSERT INTO accounts (user_id, provider, provider_account_id, password) VALUES ($1, $2, $3, $4) RETURNING *;
