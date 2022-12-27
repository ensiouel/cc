CREATE TABLE IF NOT EXISTS tokens (
    user_id       UUID REFERENCES users (id),
    refresh_token UUID NOT NULL UNIQUE
);