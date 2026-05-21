-- +goose Up
CREATE TABLE job_parts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  inventory_item_id UUID NOT NULL REFERENCES inventory_items(id),
  quantity_used INTEGER NOT NULL DEFAULT 1,
  unit_price_at_use DECIMAL(10, 2) NOT NULL, -- Snapshot of price at time of use
  subtotal DECIMAL(10, 2) GENERATED ALWAYS AS (quantity_used * unit_price_at_use) STORED,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
);

-- +goose Down
DROP TABLE job_parts;
