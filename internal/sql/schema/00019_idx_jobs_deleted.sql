-- +goose Up
CREATE INDEX idx_jobs_deleted ON jobs(deleted_at);

-- +goose Down
DROP INDEX idx_jobs_deleted;
