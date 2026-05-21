-- +goose Up
CREATE INDEX idx_audit_created ON audit_logs(created_at);

-- +goose Down
DROP INDEX idx_audit_created;
