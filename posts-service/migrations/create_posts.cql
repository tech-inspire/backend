CREATE KEYSPACE IF NOT EXISTS posts
    WITH REPLICATION = {
        'class' : 'NetworkTopologyStrategy',
        'DC1' : 3
        };

CREATE TYPE IF NOT EXISTS posts.image_variant
    (
        variant_type text,
        url          text,
        width        int,
        height       int,
        size         bigint
    );

// Main posts table keyed by post_id
CREATE TABLE IF NOT EXISTS posts.posts_by_id
(
    post_id               uuid PRIMARY KEY,
    author_id             uuid,
    images                list<frozen<image_variant>>,
    soundcloud_song       text,
    soundcloud_song_start int,
    description           text,
    created_at            timestamp
);
