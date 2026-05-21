-- +goose Up
CREATE INDEX idx_audit_user ON audit_logs(user_id);

-- +goose Down
DROP INDEX idx_audit_user;
