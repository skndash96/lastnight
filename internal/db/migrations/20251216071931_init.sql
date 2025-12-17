-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS accounts (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider TEXT NOT NULL, -- e.g. local, OAuth, etc.
  provider_account_id TEXT NOT NULL, -- e.g. OAuth provider user ID
  password BYTEA, -- only for local provider, hashed
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT unique_user_provider UNIQUE (user_id, provider)
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    domain TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE TEAM_USER_ROLE AS ENUM ('member', 'mod');

CREATE TABLE IF NOT EXISTS team_memberships (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEAM_USER_ROLE NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_team_user UNIQUE (user_id) -- Ensures a user can only join a team once
);

CREATE TABLE IF NOT EXISTS tag_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    data_type TEXT NOT NULL CHECK (data_type IN ('string', 'number', 'boolean')),
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- below table is useful for predefined tag values (e.g., for dropdowns)
CREATE TABLE IF NOT EXISTS tag_values (
    id SERIAL PRIMARY KEY,
    value TEXT UNIQUE NOT NULL,
    tag_type_id INTEGER NOT NULL REFERENCES tag_types(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_member_preferences (
    id SERIAL PRIMARY KEY,
    team_membership_id INTEGER NOT NULL REFERENCES team_memberships(id) ON DELETE CASCADE,
    tag_key TEXT NOT NULL REFERENCES tag_types(name) ON DELETE CASCADE,
    tag_value TEXT NOT NULL REFERENCES tag_values(value) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT unique_team_member_preference UNIQUE (team_membership_id, tag_key)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_member_preferences;
DROP TABLE IF EXISTS tag_values;
DROP TABLE IF EXISTS tag_types;
DROP TABLE IF EXISTS team_memberships;
DROP TABLE IF EXISTS teams;
DROP TYPE IF EXISTS TEAM_ROLE;

DROP INDEX IF EXISTS idx_accounts_user_id;

DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
