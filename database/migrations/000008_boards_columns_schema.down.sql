DROP INDEX IF EXISTS idx_time_entries_user;
DROP INDEX IF EXISTS idx_time_entries_task;
DROP TABLE IF EXISTS time_entries;

DROP INDEX IF EXISTS idx_checklist_items_list;
DROP TABLE IF EXISTS checklist_items;

DROP INDEX IF EXISTS idx_checklists_task;
DROP TABLE IF EXISTS checklists;

DROP TABLE IF EXISTS card_labels;

DROP INDEX IF EXISTS uq_labels_name_per_board;
DROP TABLE IF EXISTS labels;

DROP INDEX IF EXISTS idx_tasks_position;
DROP INDEX IF EXISTS idx_tasks_board;
DROP INDEX IF EXISTS idx_tasks_column;
ALTER TABLE tasks
    DROP COLUMN IF EXISTS remaining_estimate_seconds,
    DROP COLUMN IF EXISTS original_estimate_seconds,
    DROP COLUMN IF EXISTS story_points,
    DROP COLUMN IF EXISTS position,
    DROP COLUMN IF EXISTS column_id,
    DROP COLUMN IF EXISTS board_id;

DROP INDEX IF EXISTS idx_columns_board;
DROP INDEX IF EXISTS uq_columns_name_per_board;
DROP TABLE IF EXISTS columns;

DROP INDEX IF EXISTS idx_boards_project;
DROP TABLE IF EXISTS boards;
