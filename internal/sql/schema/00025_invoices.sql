-- +goose Up
CREATE TABLE invoices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  invoice_number VARCHAR(50) UNIQUE NOT NULL, -- Human-readable like "INV-2026-001"
  job_id UUID NOT NULL REFERENCES jobs(id),
  customer_id UUID NOT NULL REFERENCES customers(id),
  subtotal DECIMAL(10, 2) NOT NULL,
  tax_rate DECIMAL(5, 4) DEFAULT 0.05, -- GST 5%
  tax_amount DECIMAL(10, 2) GENERATED ALWAYS AS (subtotal * tax_rate) STORED,
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, paid, overdue, void
  due_date DATE NOT NULL,
  issued_date DATE NOT NULL DEFAULT CURRENT_DATE,
  paid_date DATE,
  payment_method VARCHAR(50), -- cash, credit_card, cheque, e-transfer
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE invoices;
