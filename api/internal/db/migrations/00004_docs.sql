-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION vector;

CREATE TYPE doc_proc_status AS ENUM ('pending', 'completed', 'failed');

CREATE TABLE IF NOT EXISTS docs (
  id SERIAL PRIMARY KEY,
  storage_key TEXT UNIQUE NOT NULL,
  file_sha256 TEXT NOT NULL,
  file_size BIGINT NOT NULL,
  file_mime_type TEXT NOT NULL,
  status doc_proc_status NOT NULL DEFAULT 'pending',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (file_sha256, file_size)
);

CREATE TABLE IF NOT EXISTS doc_refs (
  id SERIAL PRIMARY KEY,
  doc_id INTEGER NOT NULL REFERENCES docs(id),
  user_id INTEGER NOT NULL REFERENCES users(id),
  team_id INTEGER NOT NULL REFERENCES teams(id),
  file_name TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS doc_ref_tags (
  id SERIAL PRIMARY KEY,
  doc_ref_id INTEGER NOT NULL REFERENCES doc_refs(id),
  key_id INTEGER NOT NULL REFERENCES tag_keys(id),
  value_id INTEGER NOT NULL REFERENCES tag_values(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  UNIQUE (doc_ref_id, key_id)
);

CREATE TABLE IF NOT EXISTS doc_chunks (
  id SERIAL PRIMARY KEY,
  doc_id INTEGER NOT NULL REFERENCES doc_refs(id),
  chunk_idx INTEGER NOT NULL,
  content TEXT NOT NULL,
  meta JSONB NOT NULL DEFAULT '{}',
  embedding vector(384)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS doc_embeddings;
DROP TABLE IF EXISTS doc_ref_tags;
DROP TABLE IF EXISTS doc_refs;
DROP TABLE IF EXISTS docs;

DROP TYPE IF EXISTS doc_proc_status;

DROP EXTENSION IF EXISTS vector;
-- +goose StatementEnd
