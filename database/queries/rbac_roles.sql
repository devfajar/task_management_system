-- name: CreateRole :one
INSERT INTO roles (key, name)
VALUES ($1, $2)
RETURNING id, key, name, created_at;

-- name: GetRoleByKey :one
SELECT id, key, name, created_at FROM roles WHERE key = $1 LIMIT 1;

-- name: ListRoles :many
SELECT id, key, name, created_at FROM roles ORDER BY created_at DESC;

-- name: DeleteRole :exec
DELETE FROM roles WHERE id = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RevokeRoleFromUser :exec
DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2;

-- name: ListUserRoles :many
SELECT r.id, r.key, r.name, r.created_at
FROM user_roles ur
         JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY r.created_at DESC;
