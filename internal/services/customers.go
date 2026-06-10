// internal/services/customers.go
package services

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ValidateCompanyName(companyName string) pgtype.Text {
	cn := pgtype.Text{
		String: companyName,
		Valid:  true,
	}
	if companyName == "" {
		cn.Valid = false
	}
	return cn
}

func ValidateAddressLine2(addLine2 string) pgtype.Text {
	al2 := pgtype.Text{
		String: addLine2,
		Valid:  true,
	}
	if addLine2 == "" {
		al2.Valid = false
	}
	return al2
}

func ValidateCountry(country string) pgtype.Text {
	c := pgtype.Text{
		String: country,
		Valid:  true,
	}
	if country == "" {
		c.Valid = false
	}
	return c
}

func ValidateUserID(userID *uuid.UUID) string {
	if userID == nil {
		return ""
	}
	return userID.String()
}
