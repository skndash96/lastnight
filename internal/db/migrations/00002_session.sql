-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id INTEGER NOT NULL REFERENCES users(id),
	email TEXT NOT NULL,
	expiry TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
