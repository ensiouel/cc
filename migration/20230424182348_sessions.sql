-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sessions
(
    id            UUID DEFAULT GEN_RANDOM_UUID() NOT NULL PRIMARY KEY,
    user_id       UUID REFERENCES users (id) ON DELETE CASCADE,
    refresh_token UUID                           NOT NULL,
    ip            TEXT                           NOT NULL,
    created_at    TIMESTAMPTZ                    NOT NULL,
    updated_at    TIMESTAMPTZ                    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions CASCADE;
-- +goose StatementEnd
