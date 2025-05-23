-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts_search_info
(
    post_id         UUID PRIMARY KEY,
    author_id       UUID      NOT NULL,

    description     TEXT      NOT NULL,

    image_path      TEXT      NOT NULL,
    image_width     INT       NOT NULL,
    image_height    INT       NOT NULL,
    image_embedding vector(512),

    ratio           NUMERIC GENERATED ALWAYS AS (
        image_width::NUMERIC / NULLIF(image_height, 0)
        ) STORED,


    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_posts_image_embedding_flat
    ON posts_search_info
        USING hnsw (image_embedding vector_l2_ops);
CREATE INDEX IF NOT EXISTS idx_posts_ratio ON posts_search_info (ratio);
CREATE INDEX IF NOT EXISTS idx_posts_author_id
    ON posts_search_info (author_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_posts_ratio;
DROP INDEX idx_posts_author_id;
DROP INDEX idx_posts_image_embedding_flat;
DROP TABLE IF EXISTS posts_search_info;
-- +goose StatementEnd
