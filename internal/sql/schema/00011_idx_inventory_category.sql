-- +goose Up
CREATE INDEX idx_inventory_category ON inventory_items(category);

-- +goose Down
DROP INDEX idx_inventory_category;
