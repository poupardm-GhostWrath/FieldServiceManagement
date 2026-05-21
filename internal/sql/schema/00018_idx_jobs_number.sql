-- +goose Up
CREATE INDEX idx_jobs_number ON jobs(job_number);

-- +goose Down
DROP INDEX idx_jobs_number;
