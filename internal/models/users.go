// internal/models/users.go
package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRole struct {
	UserID   uuid.UUID `json:"user_id"`
	RoleID   int32     `json:"role_id"`
	RoleName string    `json:"role_name"`
}

type User struct {
	ID        string             `json:"id"`
	Email     string             `json:"email"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Phone     pgtype.Text        `json:"phone"`
	IsActive  pgtype.Bool        `json:"is_active"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
	Roles     []string           `json:"roles"`
}
