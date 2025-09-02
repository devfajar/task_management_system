-- name: AddTimeEntry :one
INSERT INTO time_entries (task_id, user_id, spent_seconds, work_date, note)
VALUES ($1, $2, $3, COALESCE($4, CURRENT_DATE), COALESCE($5, ''))
RETURNING id, task_id, user_id, spent_seconds, work_date, note, created_at;

-- name: SumTimeByTask :one
SELECT COALESCE(SUM(spent_seconds), 0) AS total_spent_seconds
FROM time_entries
WHERE task_id = $1;

-- name: ListTimeEntriesByTask :many
SELECT id, task_id, user_id, spent_seconds, work_date, note, created_at
FROM time_entries
WHERE task_id = $1
ORDER BY work_date DESC, created_at DESC
LIMIT $2 OFFSET $3;
