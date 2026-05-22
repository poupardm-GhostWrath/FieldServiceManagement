// internal/models/api.go
package models

import (
	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
)

type APIConfig struct {
	User      UserRole
	DB        *pgx.Conn
	DBQueries *database.Queries
	Platform  string
}
