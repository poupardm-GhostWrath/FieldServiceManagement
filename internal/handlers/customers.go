// internal/handlers/customers.go
package handlers

import (
	"net/http"

	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Customers []models.Customer `json:"customers"`
	}

	// 1. Get Customers from DB
	dbCustomers, err := config.APICfg.DBQueries.GetCustomers(r.Context())
	if err != nil {
		http.Error(w, "Couldn't retrieve customers", http.StatusInternalServerError)
		return
	}
	if len(dbCustomers) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 2. Craft payload
	customers := []models.Customer{}
	for _, dbCustomer := range dbCustomers {
		customers = append(customers, models.Customer{
			ID:           dbCustomer.ID,
			CompanyName:  dbCustomer.CompanyName,
			ContactName:  dbCustomer.ContactName,
			Email:        dbCustomer.Email,
			Phone:        dbCustomer.Phone,
			AddressLine1: dbCustomer.AddressLine1,
			AddressLine2: dbCustomer.AddressLine2,
			City:         dbCustomer.City,
			Province:     dbCustomer.Province,
			PostalCode:   dbCustomer.PostalCode,
			Country:      dbCustomer.Country,
			Notes:        dbCustomer.Notes,
			UserID:       dbCustomer.UserID.String(),
			CreatedAt:    dbCustomer.CreatedAt,
			UpdatedAt:    dbCustomer.UpdatedAt,
			DeletedAt:    dbCustomer.DeletedAt,
		})
	}

	// 3. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Customers: customers,
	})
}
