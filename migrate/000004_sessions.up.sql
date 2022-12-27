CREATE TABLE IF NOT EXISTS sessions (
    id            UUID        DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id       UUID        REFERENCES users (id) ON DELETE CASCADE,
    refresh_token UUID        NOT NULL,
    ip            TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL,
    updated_at    TIMESTAMPTZ NOT NULL
);