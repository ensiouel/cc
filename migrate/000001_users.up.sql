CREATE TABLE IF NOT EXISTS users (
    id       UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    name     TEXT  UNIQUE      NOT NULL,
    email    TEXT              NOT NULL,
    password BYTEA             NOT NULL
);