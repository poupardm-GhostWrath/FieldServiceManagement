-- +goose Up
CREATE INDEX idx_invoices_number ON invoices(invoice_number);

-- +goose Down
DROP INDEX idx_invoices_number;
