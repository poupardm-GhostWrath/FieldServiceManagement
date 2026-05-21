-- +goose Up
CREATE TABLE jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  job_number VARCHAR(50) UNIQUE NOT NULL, -- Human-readable "JOB-2026-001"
  customer_id UUID NOT NULL REFERENCES customers(id),
  assigned_technician_id UUID REFERENCES users(id),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status VARCHAR(50) NOT NULL DEFAULT 'scheduled', -- scheduled, dispatched, arrived, in_progress, completed, cancelled
  priority VARCHAR(20) DEFAULT 'normal', -- low, normal, high, emergency
  scheduled_start TIMESTAMP WITH TIME ZONE NOT NULL,
  scheduled_end TIMESTAMP WITH TIME ZONE,
  actual_start TIMESTAMP WITH TIME ZONE,
  actual_end TIMESTAMP WITH TIME ZONE,
  service_address_line1 VARCHAR(255) NOT NULL,
  service_address_line2 VARCHAR(255),
  service_city VARCHAR(100) NOT NULL,
  service_province VARCHAR(100) NOT NULL,
  service_postal_code VARCHAR(20) NOT NULL,
  notes TEXT,
  internal_notes TEXT, -- Hidden from customers
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE -- Soft delete
);

-- +goose Down
DROP TABLE jobs;
