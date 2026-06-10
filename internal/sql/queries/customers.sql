-- name: GetCustomers :many
SELECT * FROM customers;

-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1;

-- name: CreateCustomer :one
INSERT INTO customers (
  company_name, 
  contact_name, 
  email, 
  phone, 
  address_line1, 
  address_line2,
  city,
  province,
  postal_code,
  country)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET company_name = $2, 
    contact_name = $3, 
    email = $4, 
    phone = $5, 
    address_line1 = $6, 
    address_line2 = $7, 
    city = $8, 
    province = $9, 
    postal_code = $10,
    country = $11,
    notes = $12,
    user_id = $13,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
UPDATE customers
SET updated_at = NOW(), deleted_at = NOW()
WHERE id = $1;
