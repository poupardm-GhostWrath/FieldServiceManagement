-- +goose Up
CREATE TABLE roles (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL, -- 'admin', 'dispatcher', 'technician', 'customer'
  description TEXT
);

-- +goose Down
DROP TABLE roles;
