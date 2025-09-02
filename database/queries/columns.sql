-- name: CreateColumn :one
INSERT INTO columns (board_id, name, wip_limit, position)
VALUES (sqlc.arg(board_id)::uuid,
        sqlc.arg(name),
        sqlc.narg(wip_limit)::int4, -- ← nullable INT, cast ke int4
        sqlc.arg(position)::float8 -- ← NOT NULL double precision
       )
RETURNING id, board_id, name, wip_limit, position, created_at, updated_at;

-- name: ListColumnsByBoard :many
SELECT id, board_id, name, wip_limit, position, created_at, updated_at
FROM columns
WHERE board_id = $1
  AND deleted_at IS NULL
ORDER BY position ASC;

-- name: UpdateColumn :one
UPDATE columns
SET name       = COALESCE(sqlc.narg(name), name),
    wip_limit  = COALESCE(sqlc.narg(wip_limit)::int4, wip_limit), -- ← cast
    position   = COALESCE(sqlc.narg(position)::float8, position), -- ← cast
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, board_id, name, wip_limit, position, created_at, updated_at;

-- name: SoftDeleteColumn :exec
UPDATE columns
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: NextColumnPosition :one
SELECT f_next_column_position($1) AS next_pos;
