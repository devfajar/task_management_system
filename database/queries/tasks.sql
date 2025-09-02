-- name: CreateTask :one
INSERT INTO tasks (project_id, title, description, status, priority, assignee_id, due_date)
VALUES ($1, $2, COALESCE($3, ''), COALESCE($4, 'new'), COALESCE($5, 0), $6, $7)
RETURNING id, project_id, title, description, status, priority, assignee_id, due_date, created_at, updated_at;

-- name: GetTaskByID :one
SELECT id,
       project_id,
       title,
       description,
       status,
       priority,
       assignee_id,
       due_date,
       created_at,
       updated_at
FROM tasks
WHERE id = $1
  AND deleted_at IS NULL
LIMIT 1;

-- name: ListTasks :many
SELECT id,
       project_id,
       title,
       description,
       status,
       priority,
       assignee_id,
       due_date,
       created_at,
       updated_at
FROM tasks
WHERE deleted_at IS NULL
  AND (sqlc.narg(project_id)::uuid IS NULL OR project_id = sqlc.narg(project_id)::uuid)
  AND (sqlc.narg(status)::text IS NULL OR status = sqlc.narg(status)::text)
  AND (sqlc.narg(assignee_id)::uuid IS NULL OR assignee_id = sqlc.narg(assignee_id)::uuid)
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateTask :one
UPDATE tasks
SET title       = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    status      = COALESCE(sqlc.narg(status), status),
    priority    = COALESCE(sqlc.narg(priority), priority),
    assignee_id = COALESCE(sqlc.narg(assignee_id), assignee_id),
    due_date    = COALESCE(sqlc.narg(due_date), due_date),
    updated_at  = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING id, project_id, title, description, status, priority, assignee_id, due_date, created_at, updated_at;

-- name: SoftDeleteTask :exec
UPDATE tasks
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ForceDeleteTask :exec
DELETE
FROM tasks
WHERE id = $1;

-- name: ListTasksByProject :many
SELECT id,
       project_id,
       title,
       description,
       status,
       priority,
       assignee_id,
       due_date,
       created_at,
       updated_at
FROM tasks
WHERE project_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: AssignTask :exec
UPDATE tasks
SET assignee_id = $2,
    updated_at  = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ChangeTaskStatus :exec
UPDATE tasks
SET status     = $2,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;
