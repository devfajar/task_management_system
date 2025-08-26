-- users table
CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    username   TEXT        NOT NULL,
    email      TEXT        NOT NULL,
    password   TEXT        NOT NULL,
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL     DEFAULT NULL
);

-- case-insensitive uniqueness
CREATE UNIQUE INDEX users_email_unique ON users ((LOWER(email)));