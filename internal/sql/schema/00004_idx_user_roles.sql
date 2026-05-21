-- +goose Up
CREATE INDEX idx_user_roles_user ON user_roles(user_id);

-- +goose Down
DROP INDEX idx_user_roles_user;
