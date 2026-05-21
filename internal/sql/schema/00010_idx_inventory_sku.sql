-- +goose Up
CREATE INDEX idx_inventory_sku ON inventory_items(sku);

-- +goose Down
DROP INDEX idx_inventory_sku;
