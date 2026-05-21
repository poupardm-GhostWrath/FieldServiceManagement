-- +goose Up
CREATE INDEX idx_jobs_scheduled ON jobs(scheduled_start);

-- +goose Down
DROP INDEX idx_jobs_scheduled;
