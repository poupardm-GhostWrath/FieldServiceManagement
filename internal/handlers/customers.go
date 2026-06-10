// internal/handlers/customers.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/services"
)

type customerParams struct {
	CompanyName  pgtype.Text `json:"company_name"`
	ContactName  string      `json:"contact_name"`
	Email        string      `json:"email"`
	Phone        string      `json:"phone"`
	AddressLine1 string      `json:"address_line_1"`
	AddressLine2 pgtype.Text `json:"address_line_2"`
	City         string      `json:"city"`
	Province     string      `json:"province"`
	PostalCode   string      `json:"postal_code"`
	Country      pgtype.Text `json:"country"`
	Notes        pgtype.Text `json:"notes"`
	UserID       *uuid.UUID  `json:"user_id"`
}

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
		if dbCustomer.DeletedAt.Valid {
			continue // Skip if deleted
		}
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
			UserID:       dbCustomer.UserID,
			CreatedAt:    dbCustomer.CreatedAt,
			UpdatedAt:    dbCustomer.UpdatedAt,
			DeletedAt:    dbCustomer.DeletedAt,
		})
	}

	if len(customers) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 3. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Customers: customers,
	})
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Customer models.Customer `json:"customer"`
	}

	// 1. Decode parameters from Request
	decoder := json.NewDecoder(r.Body)
	params := customerParams{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 2. Validate Email
	ok := services.ValidateEmail(params.Email)
	if !ok {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	/*
		// 3. Verify CompanyName Set
		companyName := services.ValidateCompanyName(params.CompanyName)

		// 4. Verify AddressLine2 Set
		addressLine2 := services.ValidateAddressLine2(params.AddressLine2)

		// 5. Verify Country Set
		country := services.ValidateCountry(params.Country)
	*/
	// 6. Create DB Customer
	dbCustomer, err := config.APICfg.DBQueries.CreateCustomer(r.Context(), database.CreateCustomerParams{
		CompanyName:  params.CompanyName,
		ContactName:  params.ContactName,
		Email:        params.Email,
		Phone:        params.Phone,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		City:         params.City,
		Province:     params.Province,
		PostalCode:   params.PostalCode,
		Country:      params.Country,
	})
	if err != nil {
		http.Error(w, "Couldn't create customer", http.StatusInternalServerError)
		return
	}

	// 8. Respond
	RespondWithJSON(w, http.StatusCreated, response{
		Customer: models.Customer{
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
			UserID:       dbCustomer.UserID,
			CreatedAt:    dbCustomer.CreatedAt,
			UpdatedAt:    dbCustomer.UpdatedAt,
			DeletedAt:    dbCustomer.DeletedAt,
		},
	})
}

func GetCustomer(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Customer models.Customer `json:"customer"`
		// ADD: Jobs
	}

	// 1. Fetch Customer ID
	customerIDString := r.PathValue("customerID")
	customerID, err := uuid.Parse(customerIDString)
	if err != nil {
		http.Error(w, "Couldn't parse customer ID", http.StatusBadRequest)
		return
	}

	// 2. Get Customer From DB
	dbCustomer, err := config.APICfg.DBQueries.GetCustomer(r.Context(), customerID)
	if err != nil {
		http.Error(w, "Couldn't get customer", http.StatusBadRequest)
		return
	}

	// 3. Check if customer is deleted
	if dbCustomer.DeletedAt.Valid {
		http.Error(w, "Invalid customer", http.StatusBadRequest)
		return
	}

	// 5. Create payload
	customer := models.Customer{
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
		UserID:       dbCustomer.UserID,
		CreatedAt:    dbCustomer.CreatedAt,
		UpdatedAt:    dbCustomer.UpdatedAt,
		DeletedAt:    dbCustomer.DeletedAt,
	}

	// ADD: Jobs History

	// 6. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Customer: customer,
	})
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch Customer ID
	customerIDString := r.PathValue("customerID")
	customerID, err := uuid.Parse(customerIDString)
	if err != nil {
		http.Error(w, "Couldn't parse customer ID", http.StatusBadRequest)
		return
	}

	// 2. Delete Customer
	err = config.APICfg.DBQueries.DeleteCustomer(r.Context(), customerID)
	if err != nil {
		http.Error(w, "Couldn't delete customer", http.StatusInternalServerError)
		return
	}

	// 3. Respond
	w.WriteHeader(http.StatusNoContent)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Customer models.Customer `json:"customer"`
	}
	// 1. Fetch Customer ID
	customerIDString := r.PathValue("customerID")
	customerID, err := uuid.Parse(customerIDString)
	if err != nil {
		http.Error(w, "Couldn't retrieve customer ID", http.StatusBadRequest)
		return
	}

	// 2. Decode Request
	decoder := json.NewDecoder(r.Body)
	params := customerParams{}
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 3. Verify data
	ok := services.ValidateEmail(params.Email)
	if !ok {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// 4. Update Customer
	dbCustomer, err := config.APICfg.DBQueries.UpdateCustomer(r.Context(), database.UpdateCustomerParams{
		ID:           customerID,
		CompanyName:  params.CompanyName,
		ContactName:  params.ContactName,
		Email:        params.Email,
		Phone:        params.Phone,
		AddressLine1: params.AddressLine1,
		AddressLine2: params.AddressLine2,
		City:         params.City,
		Province:     params.Province,
		PostalCode:   params.PostalCode,
		Country:      params.Country,
		Notes:        params.Notes,
		UserID:       params.UserID,
	})
	if err != nil {
		http.Error(w, "Couldn't update customer", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Customer: models.Customer{
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
			UserID:       dbCustomer.UserID,
			CreatedAt:    dbCustomer.CreatedAt,
			UpdatedAt:    dbCustomer.UpdatedAt,
			DeletedAt:    dbCustomer.DeletedAt,
		},
	})
}
