-- name: CreateLabel :one
INSERT INTO labels (board_id, name, color)
VALUES ($1, $2, COALESCE($3, '#DADADA'))
RETURNING id, board_id, name, color, created_at;

-- name: ListLabelsByBoard :many
SELECT id, board_id, name, color, created_at
FROM labels
WHERE board_id = $1
ORDER BY name ASC;

-- name: AddLabelToTask :exec
INSERT INTO card_labels (task_id, label_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveLabelFromTask :exec
DELETE
FROM card_labels
WHERE task_id = $1
  AND label_id = $2;

-- name: ListTaskLabels :many
SELECT l.id, l.board_id, l.name, l.color, l.created_at
FROM card_labels cl
         JOIN labels l ON l.id = cl.label_id
WHERE cl.task_id = $1
ORDER BY l.name ASC;
