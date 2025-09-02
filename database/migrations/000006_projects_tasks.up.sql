-- pastikan pgcrypto ada (untuk gen_random_uuid v4)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- projects
CREATE TABLE projects
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    owner_id    UUID        NULL REFERENCES users (id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ NULL
);

-- tasks
CREATE TABLE tasks
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    project_id  UUID        NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    status      TEXT        NOT NULL DEFAULT 'new', -- pakai CHECK agar fleksibel tanpa enum type
    priority    INT         NOT NULL DEFAULT 0,
    assignee_id UUID        NULL REFERENCES users (id) ON DELETE SET NULL,
    due_date    TIMESTAMPTZ NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ NULL,
    CONSTRAINT tasks_status_check CHECK (status IN ('new', 'in_progress', 'done', 'blocked'))
);

-- Indeks yang sering dipakai
CREATE INDEX IF NOT EXISTS idx_projects_owner ON projects (owner_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks (project_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks (assignee_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_status_due ON tasks (status, due_date) WHERE deleted_at IS NULL;
