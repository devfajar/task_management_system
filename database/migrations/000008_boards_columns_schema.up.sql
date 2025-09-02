-- butuh pgcrypto (sudah di-enable sebelumnya)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1) Boards (per project)
CREATE TABLE boards
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    project_id  UUID        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ NULL
);
CREATE INDEX IF NOT EXISTS idx_boards_project ON boards (project_id) WHERE deleted_at IS NULL;

-- 2) Columns/Lists (per board) + ordering
CREATE TABLE columns
(
    id         UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    board_id   UUID             NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    name       TEXT             NOT NULL,
    wip_limit  INT              NULL CHECK (wip_limit IS NULL OR wip_limit >= 0),
    position   DOUBLE PRECISION NOT NULL, -- ordering key (fractional)
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ      NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_columns_name_per_board
    ON columns (board_id, LOWER(name)) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_columns_board ON columns (board_id) WHERE deleted_at IS NULL;

-- 3) Tasks upgrade: link to board/column, ordering, story points & estimates
ALTER TABLE tasks
    ADD COLUMN IF NOT EXISTS board_id                   UUID             NULL REFERENCES boards (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS column_id                  UUID             NULL REFERENCES columns (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS position                   DOUBLE PRECISION NULL,
    ADD COLUMN IF NOT EXISTS story_points               DOUBLE PRECISION NULL CHECK (story_points IS NULL OR story_points >= 0),
    ADD COLUMN IF NOT EXISTS original_estimate_seconds  INTEGER          NOT NULL DEFAULT 0 CHECK (original_estimate_seconds >= 0),
    ADD COLUMN IF NOT EXISTS remaining_estimate_seconds INTEGER          NOT NULL DEFAULT 0 CHECK (remaining_estimate_seconds >= 0);

CREATE INDEX IF NOT EXISTS idx_tasks_column ON tasks (column_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_board ON tasks (board_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_position ON tasks (column_id, position) WHERE deleted_at IS NULL;

-- 4) Labels (per board) & card_labels (many-to-many)
CREATE TABLE labels
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    board_id   UUID        NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    name       TEXT        NOT NULL,
    color      TEXT        NOT NULL DEFAULT '#DADADA',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS uq_labels_name_per_board
    ON labels (board_id, LOWER(name));

CREATE TABLE card_labels
(
    task_id  UUID NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    label_id UUID NOT NULL REFERENCES labels (id) ON DELETE CASCADE,
    PRIMARY KEY (task_id, label_id)
);

-- 5) Checklists & items
CREATE TABLE checklists
(
    id         UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    task_id    UUID             NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    title      TEXT             NOT NULL,
    position   DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checklists_task ON checklists (task_id);

CREATE TABLE checklist_items
(
    id           UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    checklist_id UUID             NOT NULL REFERENCES checklists (id) ON DELETE CASCADE,
    content      TEXT             NOT NULL,
    is_done      BOOLEAN          NOT NULL DEFAULT FALSE,
    position     DOUBLE PRECISION NOT NULL,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_checklist_items_list ON checklist_items (checklist_id);

-- 6) Time entries (time tracking)
CREATE TABLE time_entries
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    task_id       UUID        NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    user_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    spent_seconds INTEGER     NOT NULL CHECK (spent_seconds > 0),
    work_date     DATE        NOT NULL DEFAULT CURRENT_DATE,
    note          TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_time_entries_task ON time_entries (task_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_user ON time_entries (user_id);
