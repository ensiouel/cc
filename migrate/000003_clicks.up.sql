CREATE TABLE IF NOT EXISTS clicks (
    shorten_id BIGINT REFERENCES shortens (id) ON DELETE CASCADE,
    platform   TEXT,
    os         TEXT,
    referrer   TEXT,
    timestamp  TIMESTAMPTZ NOT NULL
);