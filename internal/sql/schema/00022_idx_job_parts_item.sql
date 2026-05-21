-- +goose Up
CREATE INDEX idx_job_parts_item ON job_parts(inventory_item_id);

-- +goose Down
DROP INDEX idx_job_parts_item;
