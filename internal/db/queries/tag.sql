-- name: ListTags :many
SELECT * FROM tags WHERE team_id = $1;

-- name: ListTagValues :many
SELECT * FROM tag_values WHERE tag_id = $1;

-- name: CreateTag :one
INSERT INTO tags (team_id, name, data_type) VALUES ($1, $2, $3) RETURNING *;

-- name: CreateTagValue :one
INSERT INTO tag_values (tag_id, value) VALUES ($1, $2) RETURNING *;

-- name: UpdateTag :one
UPDATE tags SET name = $2, data_type = $3 WHERE id = $1 RETURNING *;

-- name: DeleteTag :one
DELETE FROM tags WHERE id = $1 RETURNING *;

-- name: DeleteTagValue :one
DELETE FROM tag_values WHERE id = $1 RETURNING *;
