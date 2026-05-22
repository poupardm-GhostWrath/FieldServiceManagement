-- name: CreateUser :one
INSERT INTO users (email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateUserWithPhone :one
INSERT INTO users (email, password_hash, first_name, last_name, phone)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
