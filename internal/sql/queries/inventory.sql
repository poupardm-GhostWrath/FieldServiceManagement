-- name: GetInventory :many
SELECT * FROM inventory_items;

-- name: GetInventoryItemBySKU :one
SELECT * FROM inventory_items
WHERE sku = $1;

-- name: CreateInventoryItem :one
INSERT INTO inventory_items (sku, name, description, category, unit_price, quantity_in_stock, reorder_threshold, supplier_name, supplier_contact)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateInventoryItem :one
UPDATE inventory_items
SET name = $2, 
    description = $3, 
    category = $4, 
    unit_price = $5, 
    quantity_in_stock = $6, 
    reorder_threshold = $7, 
    supplier_name = $8, 
    supplier_contact = $9, 
    is_active = $10,
    updated_at = NOW()
WHERE sku = $1
RETURNING *;
