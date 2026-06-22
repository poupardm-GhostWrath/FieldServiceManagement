// internal/handlers/inventory.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

type inventoryItemParams struct {
	SKU              string         `json:"sku"`
	Name             string         `json:"name"`
	Description      pgtype.Text    `json:"description"`
	Category         pgtype.Text    `json:"category"`
	UnitPrice        pgtype.Numeric `json:"unit_price"`
	QuantityInStock  int32          `json:"quantity_in_stock"`
	ReorderThreshold int32          `json:"reorder_threshold"`
	SupplierName     pgtype.Text    `json:"supplier_name"`
	SupplierContact  pgtype.Text    `json:"supplier_contact"`
	IsActive         pgtype.Bool    `json:"is_active"`
}

func CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	type response struct {
		InventoryItem models.InventoryItem
	}
	// 1. Decode parameters
	decoder := json.NewDecoder(r.Body)
	params := inventoryItemParams{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 2. Create Inventory Item
	dbInventoryItem, err := config.APICfg.DBQueries.CreateInventoryItem(r.Context(), database.CreateInventoryItemParams{
		Sku:              params.SKU,
		Name:             params.Name,
		Description:      params.Description,
		Category:         params.Category,
		UnitPrice:        params.UnitPrice,
		QuantityInStock:  params.QuantityInStock,
		ReorderThreshold: params.ReorderThreshold,
		SupplierName:     params.SupplierName,
		SupplierContact:  params.SupplierContact,
	})
	if err != nil {
		http.Error(w, "Couldn't create inventory item", http.StatusInternalServerError)
		return
	}

	// 3. Respond
	RespondWithJSON(w, http.StatusCreated, response{
		InventoryItem: models.InventoryItem{
			ID:               dbInventoryItem.ID,
			SKU:              dbInventoryItem.Sku,
			Name:             dbInventoryItem.Name,
			Description:      dbInventoryItem.Description,
			Category:         dbInventoryItem.Category,
			UnitPrice:        dbInventoryItem.UnitPrice,
			QuantityInStock:  dbInventoryItem.QuantityInStock,
			ReorderThreshold: dbInventoryItem.ReorderThreshold,
			SupplierName:     dbInventoryItem.SupplierName,
			SupplierContact:  dbInventoryItem.SupplierContact,
			IsActive:         dbInventoryItem.IsActive,
			CreatedAt:        dbInventoryItem.CreatedAt,
			UpdatedAt:        dbInventoryItem.UpdatedAt,
		},
	})
}

func GetInventory(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Inventory []models.InventoryItem `json:"inventory"`
	}

	// 1. Get inventory from DB
	dbInventory, err := config.APICfg.DBQueries.GetInventory(r.Context())
	if err != nil {
		http.Error(w, "Couldn't retrieve inventory", http.StatusInternalServerError)
		return
	}

	// 2. Check if empty
	if len(dbInventory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 3. Create payload
	inventory := []models.InventoryItem{}
	for _, dbItem := range dbInventory {
		if !dbItem.IsActive.Bool {
			continue
		}
		item := models.InventoryItem{
			ID:               dbItem.ID,
			SKU:              dbItem.Sku,
			Name:             dbItem.Name,
			Description:      dbItem.Description,
			Category:         dbItem.Category,
			UnitPrice:        dbItem.UnitPrice,
			QuantityInStock:  dbItem.QuantityInStock,
			ReorderThreshold: dbItem.ReorderThreshold,
			SupplierName:     dbItem.SupplierName,
			SupplierContact:  dbItem.SupplierContact,
			IsActive:         dbItem.IsActive,
			CreatedAt:        dbItem.CreatedAt,
			UpdatedAt:        dbItem.UpdatedAt,
		}
		inventory = append(inventory, item)
	}

	// 4. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Inventory: inventory,
	})
}

func GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	type response struct {
		InventoryItem models.InventoryItem
	}

	// 1. Get SKU
	sku := r.PathValue("sku")

	// 2. Get Item From DB
	dbItem, err := config.APICfg.DBQueries.GetInventoryItemBySKU(r.Context(), sku)
	if err != nil {
		http.Error(w, "Couldn't retrieve inventory item", http.StatusBadRequest)
		return
	}

	// 3. Respond
	RespondWithJSON(w, http.StatusOK, response{
		InventoryItem: models.InventoryItem{
			ID:               dbItem.ID,
			SKU:              dbItem.Sku,
			Name:             dbItem.Name,
			Description:      dbItem.Description,
			Category:         dbItem.Category,
			UnitPrice:        dbItem.UnitPrice,
			QuantityInStock:  dbItem.QuantityInStock,
			ReorderThreshold: dbItem.ReorderThreshold,
			SupplierName:     dbItem.SupplierName,
			SupplierContact:  dbItem.SupplierContact,
			IsActive:         dbItem.IsActive,
			CreatedAt:        dbItem.CreatedAt,
			UpdatedAt:        dbItem.UpdatedAt,
		},
	})
}

func UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	type response struct {
		InventoryItem models.InventoryItem
	}

	// 1. Get SKU
	sku := r.PathValue("sku")

	// 2. Decode Parameters
	decoder := json.NewDecoder(r.Body)
	params := inventoryItemParams{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 3. Get Item From DB
	dbItem, err := config.APICfg.DBQueries.GetInventoryItemBySKU(r.Context(), sku)
	if err != nil {
		http.Error(w, "Couldn't retrieve inventory item", http.StatusBadRequest)
		return
	}

	// 4. Parameters verification
	name := dbItem.Name
	if params.Name != "" {
		name = params.Name
	}

	quantityInStock := max(params.QuantityInStock, 0)

	reorderThreshold := max(params.ReorderThreshold, 0)

	// 4. Create Update payload
	updatedItem := database.UpdateInventoryItemParams{
		Sku:              dbItem.Sku,
		Name:             name,
		Description:      params.Description,
		Category:         params.Category,
		UnitPrice:        params.UnitPrice,
		QuantityInStock:  quantityInStock,
		ReorderThreshold: reorderThreshold,
		SupplierName:     params.SupplierName,
		SupplierContact:  params.SupplierContact,
		IsActive:         params.IsActive,
	}

	// 5. Update DB
	updatedDBItem, err := config.APICfg.DBQueries.UpdateInventoryItem(r.Context(), updatedItem)
	if err != nil {
		http.Error(w, "Couldn't update inventory item", http.StatusInternalServerError)
		return
	}

	// 6. Respond
	RespondWithJSON(w, http.StatusOK, response{
		InventoryItem: models.InventoryItem{
			ID:               updatedDBItem.ID,
			SKU:              updatedDBItem.Sku,
			Name:             updatedDBItem.Name,
			Description:      updatedDBItem.Description,
			Category:         updatedDBItem.Category,
			UnitPrice:        updatedDBItem.UnitPrice,
			QuantityInStock:  updatedDBItem.QuantityInStock,
			ReorderThreshold: updatedDBItem.ReorderThreshold,
			SupplierName:     updatedDBItem.SupplierName,
			SupplierContact:  updatedDBItem.SupplierContact,
			IsActive:         updatedDBItem.IsActive,
			CreatedAt:        updatedDBItem.CreatedAt,
			UpdatedAt:        updatedDBItem.UpdatedAt,
		},
	})
}
