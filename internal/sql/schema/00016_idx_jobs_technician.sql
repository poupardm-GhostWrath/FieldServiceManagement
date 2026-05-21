-- +goose Up
CREATE INDEX idx_jobs_technician ON jobs(assigned_technician_id);

-- +goose Down
DROP INDEX idx_jobs_technician;
