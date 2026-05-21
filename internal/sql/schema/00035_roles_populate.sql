-- +goose Up
INSERT INTO roles (name, description) VALUES
  ('admin', 'Full system access'),
  ('dispatcher', 'Can manage jobs and schedules'),
  ('technician', 'Can view and update assigned jobs'),
  ('customer', 'Limited access to own data');

-- +goose Down
DELETE FROM roles;
