-- name: GetUserTokenVersion :one
SELECT token_version
FROM users
WHERE id = $1
LIMIT 1;

-- name: BumpUserTokenVersion :exec
UPDATE users
SET token_version = token_version + 1
WHERE id = $1;

-- name: InsertAuditLog :one
INSERT INTO audit_logs (user_id, action, entity, entity_id, metadata, ip_address, user_agent)
VALUES (sqlc.narg(user_id)::uuid,
        sqlc.arg(action),
        sqlc.arg(entity),
        sqlc.narg(entity_id)::uuid,
        COALESCE(sqlc.narg(metadata)::jsonb, '{}'::jsonb),
        sqlc.narg(ip_address),
        sqlc.narg(user_agent))
RETURNING id, user_id, action, entity, entity_id, metadata, ip_address, user_agent, created_at;




