-- name: CreateUserRoles :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserRoles :many
SELECT roles.name FROM user_roles
INNER JOIN roles ON roles.id = role_id
WHERE user_id = $1;
