-- name: GetRoles :many
SELECT * FROM roles;

-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1;
