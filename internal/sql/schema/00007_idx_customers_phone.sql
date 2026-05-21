-- +goose Up
CREATE INDEX idx_customers_phone ON customers(phone);

-- +goose Down
DROP INDEX idx_customers_phone;
