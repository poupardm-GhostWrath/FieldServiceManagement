-- +goose Up
CREATE INDEX idx_invoices_job ON invoices(job_id);

-- +goose Down
DROP INDEX idx_invoices_job;
