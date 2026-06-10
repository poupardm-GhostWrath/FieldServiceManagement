// internal/services/phone.go
package services

import "github.com/jackc/pgx/v5/pgtype"

func ValidatePhone(phone string) pgtype.Text {
	pgPhone := pgtype.Text{
		String: phone,
		Valid:  true,
	}
	if phone == "" {
		pgPhone.Valid = false
	}
	return pgPhone
}
