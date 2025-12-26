-- name: GetOrCreateUpload :one
INSERT INTO uploads (storage_key, file_sha256, file_size, file_mime_type)
VALUES ($1, $2, $3, $4)
ON CONFLICT (file_sha256, file_size) DO UPDATE SET storage_key = uploads.storage_key
RETURNING *, (xmax = 0) as created;

-- name: GetUploadRef :one
SELECT
    upload_refs.id AS upload_ref_id,
    upload_refs.upload_id,
    upload_refs.team_id,
    upload_refs.uploader_id,
    upload_refs.file_name,
    uploads.storage_key,
    uploads.file_sha256,
    uploads.file_size,
    uploads.file_mime_type
FROM upload_refs
LEFT JOIN uploads ON upload_refs.upload_id = uploads.id
WHERE upload_refs.id = $1;

-- name: CreateUploadRef :one
INSERT INTO upload_refs (upload_id, team_id, uploader_id, file_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateUploadRefTags :exec
INSERT INTO upload_ref_tags (upload_ref_id, key_id, value_id)
VALUES ($1, $2, $3);

-- name: DeleteAllUploadRefTags :exec
DELETE FROM upload_ref_tags WHERE upload_ref_id = $1;
