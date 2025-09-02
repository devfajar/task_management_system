-- name: NextChecklistPosition :one
SELECT f_next_checklist_position($1) AS next_pos;

-- name: CreateChecklist :one
INSERT INTO checklists (task_id, title, position)
VALUES ($1, $2, $3)
RETURNING id, task_id, title, position, created_at;

-- name: ListChecklists :many
SELECT id, task_id, title, position, created_at
FROM checklists
WHERE task_id = $1
ORDER BY position ASC;

-- name: RenameChecklist :one
UPDATE checklists
SET title    = COALESCE(sqlc.narg(title), title),
    position = COALESCE(sqlc.narg(position), position)
WHERE id = $1
RETURNING id, task_id, title, position, created_at;

-- name: DeleteChecklist :exec
DELETE
FROM checklists
WHERE id = $1;

-- Items
-- name: NextChecklistItemPosition :one
SELECT f_next_checklist_item_position($1) AS next_pos;

-- name: CreateChecklistItem :one
INSERT INTO checklist_items (checklist_id, content, is_done, position)
VALUES ($1, $2, COALESCE($3, FALSE), $4)
RETURNING id, checklist_id, content, is_done, position, created_at;

-- name: UpdateChecklistItem :one
UPDATE checklist_items
SET content  = COALESCE(sqlc.narg(content), content),
    is_done  = COALESCE(sqlc.narg(is_done), is_done),
    position = COALESCE(sqlc.narg(position), position)
WHERE id = $1
RETURNING id, checklist_id, content, is_done, position, created_at;

-- name: DeleteChecklistItem :exec
DELETE
FROM checklist_items
WHERE id = $1;

-- name: ListChecklistItems :many
SELECT id, checklist_id, content, is_done, position, created_at
FROM checklist_items
WHERE checklist_id = $1
ORDER BY position ASC;
