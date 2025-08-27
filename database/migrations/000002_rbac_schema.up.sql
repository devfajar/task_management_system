-- aktifkan pgcrypto (perlu privilege; di managed DB biasanya boleh)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- roles
CREATE TABLE roles
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(), -- v4
    key        TEXT        NOT NULL UNIQUE,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- permissions
CREATE TABLE permissions
(
    id          UUID PRIMARY KEY     DEFAULT gen_random_uuid(), -- v4
    key         TEXT        NOT NULL UNIQUE,
    description TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- user_roles (many-to-many)
CREATE TABLE user_roles
(
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id    UUID        NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

-- role_permissions (many-to-many)
CREATE TABLE role_permissions
(
    role_id       UUID        NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id UUID        NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);
