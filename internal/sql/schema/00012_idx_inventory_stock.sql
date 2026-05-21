-- +goose Up
CREATE INDEX idx_inventory_stock ON inventory_items(quantity_in_stock);

-- +goose Down
DROP INDEX idx_inventory_stock;
