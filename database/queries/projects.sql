-- name: CreateProject :one
INSERT INTO projects (name, description, owner_id)
VALUES (sqlc.arg(name),
        COALESCE(sqlc.narg(description), ''),
        sqlc.narg(owner_id)::uuid)
RETURNING id, name, description, owner_id, created_at, updated_at;

-- name: GetProjectByID :one
SELECT id, name, description, owner_id, created_at, updated_at
FROM projects
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListProjects :many
SELECT id, name, description, owner_id, created_at, updated_at
FROM projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateProject :one
UPDATE projects
SET name        = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    updated_at  = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, name, description, owner_id, created_at, updated_at;

-- name: SoftDeleteProject :exec
UPDATE projects
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ForceDeleteProject :exec
DELETE
FROM projects
WHERE id = $1;
