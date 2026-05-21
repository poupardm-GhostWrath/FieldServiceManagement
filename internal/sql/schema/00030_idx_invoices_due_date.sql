-- +goose Up
CREATE INDEX idx_invoices_due_date ON invoices(due_date);

-- +goose Down
DROP INDEX idx_invoices_due_date;
