-- name: CreatePermission :one
INSERT INTO permissions (key, description)
VALUES ($1, $2)
RETURNING id, key, description, created_at;

-- name: GetPermissionByKey :one
SELECT id, key, description, created_at FROM permissions WHERE key = $1 LIMIT 1;

-- name: ListPermissions :many
SELECT id, key, description, created_at FROM permissions ORDER BY created_at DESC;

-- name: DeletePermission :exec
DELETE FROM permissions WHERE id = $1;

-- name: GrantPermissionToRole :exec
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RevokePermissionFromRole :exec
DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2;

-- name: HasPermission :one
SELECT EXISTS (
    SELECT 1
    FROM user_roles ur
             JOIN role_permissions rp ON rp.role_id = ur.role_id
             JOIN permissions p ON p.id = rp.permission_id
    WHERE ur.user_id = $1
      AND p.key = $2
) AS has_permission;

-- name: ListUserPermissionKeys :many
SELECT DISTINCT p.key
FROM user_roles ur
         JOIN role_permissions rp ON rp.role_id = ur.role_id
         JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
ORDER BY p.key ASC;
