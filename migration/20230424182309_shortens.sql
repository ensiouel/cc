-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shortens
(
    id         BIGINT PRIMARY KEY,
    url        TEXT        NOT NULL,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    title      TEXT        NOT NULL,
    tags       TEXT[]      NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (url, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shortens CASCADE;
-- +goose StatementEnd
