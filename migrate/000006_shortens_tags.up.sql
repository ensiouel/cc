CREATE TABLE IF NOT EXISTS shortens_tags
(
    shorten_id BIGINT REFERENCES shortens (id) ON DELETE CASCADE,
    tag_id     UUID REFERENCES tags (id) ON DELETE CASCADE,
    UNIQUE (shorten_id, tag_id)
);