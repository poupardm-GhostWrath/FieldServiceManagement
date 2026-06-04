// internal/config/common.go
package config

import (
	"context"
	"fmt"
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
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "fsm_user"
	}
	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "fsm_pass"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "fsm_db"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbURL := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPass, dbPort, dbName)

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
