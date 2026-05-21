-- +goose Up
CREATE TABLE job_labor (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
  technician_id UUID NOT NULL REFERENCES users(id),
  hours_worked DECIMAL(5, 2) NOT NULL, -- e.g., 2.5 hours
  hourly_rate DECIMAL(10, 2) NOT NULL,
  subtotal DECIMAL(10, 2) GENERATED ALWAYS AS (hours_worked * hourly_rate) STORED,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE job_labor;
