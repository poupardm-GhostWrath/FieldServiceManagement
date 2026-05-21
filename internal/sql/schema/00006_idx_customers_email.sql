-- +goose Up
CREATE INDEX idx_customers_email ON customers(email);

-- +goose Down
DROP INDEX idx_customers_email;
