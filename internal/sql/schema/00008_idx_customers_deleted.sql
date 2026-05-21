-- +goose Up
CREATE INDEX idx_customers_deleted ON customers(deleted_at);

-- +goose Down
DROP INDEX idx_customers_deleted;
