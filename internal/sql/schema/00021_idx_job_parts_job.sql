-- +goose Up
CREATE INDEX idx_job_parts_job ON job_parts(job_id);

-- +goose Down
DROP INDEX idx_job_parts_job;
