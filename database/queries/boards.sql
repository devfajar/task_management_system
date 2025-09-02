-- name: CreateBoard :one
INSERT INTO boards (project_id, name, description)
VALUES ($1, $2, $3)
RETURNING id, project_id, name, description, created_at, updated_at;

-- name: ListBoardsByProject :many
SELECT id, project_id, name, description, created_at, updated_at
FROM boards
WHERE project_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetBoardByID :one
SELECT id, project_id, name, description, created_at, updated_at
FROM boards
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: SoftDeleteBoard :exec
UPDATE boards SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: ForceDeleteBoard :exec
DELETE FROM boards WHERE id = $1;
