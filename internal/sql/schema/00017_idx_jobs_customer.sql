-- +goose Up
CREATE INDEX idx_jobs_customer ON jobs(customer_id);

-- +goose Down
DROP INDEX idx_jobs_customer;
