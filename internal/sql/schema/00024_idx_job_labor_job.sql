-- +goose Up
CREATE INDEX idx_job_labor_job ON job_labor(job_id);

-- +goose Down
DROP INDEX idx_job_labor_job;
