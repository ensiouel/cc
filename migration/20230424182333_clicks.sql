-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clicks
(
    shorten_id BIGINT REFERENCES shortens (id) ON DELETE CASCADE,
    platform   TEXT,
    os         TEXT,
    referer    TEXT,
    ip         TEXT,
    timestamp  TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clicks CASCADE;
-- +goose StatementEnd
