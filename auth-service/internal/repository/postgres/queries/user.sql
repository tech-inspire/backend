-- name: CreateUser :exec
INSERT INTO users (user_id, email, name, username, password_hash, description, avatar_url)
VALUES (@user_id, @email, @name, @username, @password_hash, @description, @avatar_url);

-- name: GetUserByEmail :one
SELECT sqlc.embed(users)
FROM users
WHERE email = @email;


-- name: GetUserByUsername :one
SELECT sqlc.embed(users)
FROM users
WHERE username = @username;

-- name: GetUserByID :one
SELECT sqlc.embed(users)
FROM users
WHERE users.user_id = $1;

-- name: GetUsersByIDs :many
SELECT sqlc.embed(users)
FROM users
WHERE users.user_id = ANY (@user_ids::uuid[]);


-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash,
    updated_at    = NOW()
WHERE user_id = @user_id;

-- name: DeleteUserByID :exec
DELETE
FROM users
WHERE user_id = @user_id;

-- name: UpdateUserByID :exec
UPDATE users
SET name          = COALESCE(sqlc.narg('name'), name),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    username      = COALESCE(sqlc.narg('username'), username),
    description   = COALESCE(sqlc.narg('description'), description),
    avatar_url    = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    updated_at    = NOW()
WHERE user_id = @user_id;

-- name: ClearUserAvatarURL :exec
UPDATE users SET avatar_url = NULL WHERE user_id = @user_id;