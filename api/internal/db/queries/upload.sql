-- name: GetOrCreateDoc :one
INSERT INTO docs (storage_key, file_sha256, file_size, file_mime_type)
VALUES ($1, $2, $3, $4)
ON CONFLICT (file_sha256, file_size) DO UPDATE SET storage_key = docs.storage_key
RETURNING *, (xmax = 0) as created;

-- name: GetDoc :one
SELECT *
FROM docs
WHERE docs.id = $1;

-- name: CreateDocRef :one
INSERT INTO doc_refs (doc_id, team_id, user_id, file_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateDocRefTags :exec
INSERT INTO doc_ref_tags (doc_ref_id, key_id, value_id)
VALUES ($1, $2, $3);

-- name: DeleteAllDocRefTags :exec
DELETE FROM doc_ref_tags WHERE doc_ref_id = $1;
