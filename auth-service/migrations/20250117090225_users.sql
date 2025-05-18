-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users
(
    user_id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email         VARCHAR(150) UNIQUE            NOT NULL,
    username      VARCHAR(150) UNIQUE            NOT NULL,

    name          VARCHAR(150)                   NOT NULL,
    description   VARCHAR(200)                   NOT NULL,

    avatar_url    TEXT                           NULL,

    password_hash bytea                          NOT NULL,

    is_admin      bool             DEFAULT false NOT NULL,

    created_at    TIMESTAMP        DEFAULT NOW() NOT NULL,
    updated_at    TIMESTAMP        DEFAULT NOW() NOT NULL
);

-- Create index on email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on username
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
