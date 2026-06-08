-- name: CreateUserRoles :one
WITH inserted_user_roles AS (
  INSERT INTO user_roles (user_id, role_id)
  VALUES ($1, $2)
  RETURNING *
)
SELECT inserted_user_roles.*, roles.name AS role_name
FROM inserted_user_roles
INNER JOIN roles ON roles.id = inserted_user_roles.role_id;

-- name: GetUserRoles :many
SELECT roles.name FROM user_roles
INNER JOIN roles ON roles.id = role_id
WHERE user_id = $1;

-- name: DeleteUserRoles :exec
DELETE FROM user_roles
WHERE user_id = $1;
