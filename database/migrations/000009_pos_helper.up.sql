-- Columns
CREATE OR REPLACE FUNCTION f_next_column_position(p_board_id uuid)
    RETURNS double precision
    LANGUAGE sql
    STABLE AS
$$
SELECT COALESCE(MAX(position), 0::double precision) + 1024::double precision
FROM columns
WHERE board_id = p_board_id
  AND deleted_at IS NULL
$$;

-- Tasks
CREATE OR REPLACE FUNCTION f_next_task_position(p_column_id uuid)
    RETURNS double precision
    LANGUAGE sql
    STABLE AS
$$
SELECT COALESCE(MAX(position), 0::double precision) + 1024::double precision
FROM tasks
WHERE column_id = p_column_id
  AND deleted_at IS NULL
$$;

-- Checklists
CREATE OR REPLACE FUNCTION f_next_checklist_position(p_task_id uuid)
    RETURNS double precision
    LANGUAGE sql
    STABLE AS
$$
SELECT COALESCE(MAX(position), 0::double precision) + 1024::double precision
FROM checklists
WHERE task_id = p_task_id
$$;

CREATE OR REPLACE FUNCTION f_next_checklist_item_position(p_checklist_id uuid)
    RETURNS double precision
    LANGUAGE sql
    STABLE AS
$$
SELECT COALESCE(MAX(position), 0::double precision) + 1024::double precision
FROM checklist_items
WHERE checklist_id = p_checklist_id
$$;
