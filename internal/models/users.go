// internal/models/users.go
package models

import (
	"github.com/google/uuid"
)

type UserRole struct {
	UserID   uuid.UUID `json:"user_id"`
	RoleID   int32     `json:"role_id"`
	RoleName string    `json:"role_name"`
}

type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
