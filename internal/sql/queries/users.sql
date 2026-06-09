-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUsers :many
SELECT * FROM users;

-- name: UpdateUserProfileByID :one
UPDATE users
SET email = $2, password_hash = $3, first_name = $4, last_name = $5, phone = $6, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserByID :one
UPDATE users
SET email = $2, first_name = $3, last_name = $4, phone = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteUserByID :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = $1;
