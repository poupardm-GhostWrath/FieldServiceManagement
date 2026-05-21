-- +goose Up
CREATE INDEX idx_audit_table ON audit_logs(table_name);

-- +goose Down
DROP INDEX idx_audit_table;
