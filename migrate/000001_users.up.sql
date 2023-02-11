CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    name     TEXT UNIQUE NOT NULL,
    email    TEXT        NOT NULL,
    password BYTEA       NOT NULL
);