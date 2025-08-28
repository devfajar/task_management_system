-- name: GetUserCredentialsByEmail :one
SELECT id, password, is_active, deleted_at
FROM users
WHERE LOWER(email) = LOWER($1)
LIMIT 1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at, user_agent, ip_address)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, token_hash, issued_at, expires_at, revoked_at, user_agent, ip_address;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, issued_at, expires_at, revoked_at, user_agent, ip_address
FROM refresh_tokens
WHERE token_hash = $1
LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE id = $1 AND revoked_at IS NULL;
