-- name: GetInventory :many
SELECT * FROM inventory_items;

-- name: GetInventoryItemBySKU :one
SELECT * FROM inventory_items
WHERE sku = $1;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items (sku, name, description, category, unit_price, quantity_in_stock, reorder_threshold, supplier_name, supplier_contact)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
