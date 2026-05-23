// internal/config/common.go
package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
)

var APICfg GlobalConfig

type GlobalConfig struct {
	DB        *pgx.Conn
	DBQueries *database.Queries
	Platform  string
}

func init() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "live" // 'live' if not in 'dev' mode
	}

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbQueries := database.New(db)

	APICfg = GlobalConfig{
		DB:        db,
		DBQueries: dbQueries,
		Platform:  platform,
	}
}
