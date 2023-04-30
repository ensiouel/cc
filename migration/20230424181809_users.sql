-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    name     TEXT UNIQUE NOT NULL,
    password BYTEA       NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
