-- +goose Up
CREATE INDEX idx_jobs_status ON jobs(status);

-- +goose Down
DROP INDEX idx_jobs_status;
