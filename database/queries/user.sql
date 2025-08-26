-- name: CreateUser :one
INSERT INTO users (id, username, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email, is_active, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, is_active, created_at, updated_at
FROM users WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, is_active, created_at, updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUserProfile :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email     = COALESCE(sqlc.narg(email), email),
    is_active = COALESCE(sqlc.narg(is_active), is_active),
    updated_at = NOW()
WHERE id = $1
RETURNING id, username, email, is_active, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2,
    updated_at    = NOW()
WHERE id = $1;

-- Soft delete (default)
-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- Restore (opsional)
-- name: RestoreUser :exec
UPDATE users
SET deleted_at = NULL,
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NOT NULL;
