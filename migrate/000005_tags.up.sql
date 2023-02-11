CREATE TABLE IF NOT EXISTS tags
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    title   TEXT NOT NULL,
    UNIQUE (user_id, title)
);