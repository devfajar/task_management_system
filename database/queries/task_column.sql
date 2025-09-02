-- name: NextTaskPosition :one
SELECT f_next_task_position($1) AS next_pos;

-- name: CreateTaskInColumn :one
INSERT INTO tasks (project_id, title, description, status, priority,
                   assignee_id, due_date, board_id, column_id, position,
                   story_points, original_estimate_seconds, remaining_estimate_seconds)
VALUES (sqlc.arg(project_id)::uuid,
        sqlc.arg(title),
        COALESCE(sqlc.narg(description)::text, ''),
        COALESCE(sqlc.narg(status)::text, 'new'),
        COALESCE(sqlc.narg(priority)::int4, 0),
        sqlc.narg(assignee_id)::uuid,
        sqlc.narg(due_date)::timestamptz,
        sqlc.arg(board_id)::uuid,
        sqlc.arg(column_id)::uuid,
        sqlc.arg(position)::float8,
        sqlc.narg(story_points)::float8,
        COALESCE(sqlc.narg(original_estimate_seconds)::int4, 0),
        COALESCE(sqlc.narg(remaining_estimate_seconds)::int4, 0))
RETURNING
    id, project_id, title, description, status, priority, assignee_id, due_date,
    board_id, column_id, position, story_points, original_estimate_seconds, remaining_estimate_seconds,
    created_at, updated_at;

-- name: MoveTaskToColumn :exec
UPDATE tasks
SET column_id  = $2,
    position   = $3,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ReorderTaskPosition :exec
UPDATE tasks
SET position   = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: UpdateTaskSprintFields :one
UPDATE tasks
SET title        = COALESCE(sqlc.narg(title), title),
    description  = COALESCE(sqlc.narg(description), description),
    status       = COALESCE(sqlc.narg(status), status), -- optional, bila masih dipakai
    priority     = COALESCE(sqlc.narg(priority), priority),
    assignee_id  = COALESCE(sqlc.narg(assignee_id), assignee_id),
    due_date     = COALESCE(sqlc.narg(due_date), due_date),
    story_points = COALESCE(sqlc.narg(story_points), story_points),
    updated_at   = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, project_id, title, description, status, priority, assignee_id, due_date,
    board_id, column_id, position, story_points, original_estimate_seconds, remaining_estimate_seconds,
    created_at, updated_at;

-- name: UpdateTaskEstimates :one
UPDATE tasks
SET original_estimate_seconds  = COALESCE(sqlc.narg(original_estimate_seconds), original_estimate_seconds),
    remaining_estimate_seconds = COALESCE(sqlc.narg(remaining_estimate_seconds), remaining_estimate_seconds),
    updated_at                 = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, project_id, title, description, status, priority, assignee_id, due_date,
    board_id, column_id, position, story_points, original_estimate_seconds, remaining_estimate_seconds,
    created_at, updated_at;

-- name: ListTasksByColumn :many
SELECT id,
       project_id,
       title,
       description,
       status,
       priority,
       assignee_id,
       due_date,
       board_id,
       column_id,
       position,
       story_points,
       original_estimate_seconds,
       remaining_estimate_seconds,
       created_at,
       updated_at
FROM tasks
WHERE column_id = $1
  AND deleted_at IS NULL
ORDER BY position ASC
LIMIT $2 OFFSET $3;
