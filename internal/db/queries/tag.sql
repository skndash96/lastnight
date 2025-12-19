-- name: ListTags :many
SELECT
  k.id AS key_id,
  k.name AS key,
  sv.id AS value_id,
  sv.value AS value,
  JSONB_AGG(
    JSONB_BUILD_OBJECT(
      'id', v.id,
      'value', v.value
    )
    ORDER BY v.value
  ) FILTER (WHERE v.id IS NOT NULL) AS options
  FROM tag_keys k
  LEFT JOIN team_member_tags tmt ON tmt.key_id = k.id AND tmt.membership_id = $1
  LEFT JOIN tag_values sv ON sv.id = tmt.value_id
  LEFT JOIN tag_values v ON v.key_id = k.id
  GROUP BY k.id, k.name, sv.id, sv.value;

-- name: CreateTagKey :one
INSERT INTO tag_keys (team_id, name, data_type) VALUES ($1, $2, $3) RETURNING *;

-- name: CreateTagValue :one
INSERT INTO tag_values (key_id, value) VALUES ($1, $2) RETURNING *;

-- name: UpdateTagKey :one
UPDATE tag_keys SET name = $2, data_type = $3 WHERE id = $1 RETURNING *;

-- name: DeleteTagKey :one
DELETE FROM tag_keys WHERE id = $1 RETURNING *;

-- name: DeleteTagValue :one
DELETE FROM tag_values WHERE id = $1 RETURNING *;
