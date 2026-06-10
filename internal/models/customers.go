// internal/models/customers.go
package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Customer struct {
	ID           uuid.UUID          `json:"id"`
	CompanyName  pgtype.Text        `json:"company_name"`
	ContactName  string             `json:"contact_name"`
	Email        string             `json:"email"`
	Phone        string             `json:"phone"`
	AddressLine1 string             `json:"address_line_1"`
	AddressLine2 pgtype.Text        `json:"address_line_2"`
	City         string             `json:"city"`
	Province     string             `json:"province"`
	PostalCode   string             `json:"postal_code"`
	Country      pgtype.Text        `json:"country"`
	Notes        pgtype.Text        `json:"notes"`
	UserID       *uuid.UUID         `json:"userID"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	UpdatedAt    pgtype.Timestamptz `json:"updated_at"`
	DeletedAt    pgtype.Timestamptz `json:"deleted_at"`
}
