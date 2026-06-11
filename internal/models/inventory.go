// internal/models/inventory.go
package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InventoryItem struct {
	ID               uuid.UUID          `json:"id"`
	SKU              string             `json:"sku"`
	Name             string             `json:"name"`
	Description      pgtype.Text        `json:"description"`
	Category         pgtype.Text        `json:"category"`
	UnitPrice        pgtype.Numeric     `json:"unit_price"`
	QuantityInStock  int32              `json:"quantity_in_stock"`
	ReorderThreshold int32              `json:"reorder_threshold"`
	SupplierName     pgtype.Text        `json:"supplier_name"`
	SupplierContact  pgtype.Text        `json:"supplier_contact"`
	IsActive         pgtype.Bool        `json:"is_active"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	UpdatedAt        pgtype.Timestamptz `json:"updated_at"`
}
