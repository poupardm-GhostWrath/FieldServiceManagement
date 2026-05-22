// Package main
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

func main() {
	// Get Environmental Data
	// POSTGRES Environment Variables
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		log.Fatal("POSTGRES_USER must be set")
	}
	dbPass := os.Getenv("POSTGRES_PASS")
	if dbPass == "" {
		log.Fatal("POSTGRES_PASS must be set")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		log.Fatal("POSTGRES_DB must be set")
	}
	dbAddr := os.Getenv("POSTGRES_ADDR")
	if dbAddr == "" {
		dbAddr = "localhost:5432" // Use localhost if not set
	}
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbAddr,
		dbName)

	// Go Server Environment Variables
	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if filepathRoot == "" {
		log.Fatal("FILEPATH_ROOT must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	// Development Environment Variables
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = "live" // 'live' if not in 'dev' mode
	}

	// Open SQL database
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	// Create Queries
	dbQueries := database.New(db)

	// Create APIConfig
	apiCfg := models.APIConfig{
		DB:        db,
		DBQueries: dbQueries,
		Platform:  platform,
	}

	// Create Server Mux
	mux := http.NewServeMux()

	// Set Endpoints
	appHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", appHandler)

	// Authentication & Users
	mux.HandleFunc("POST /auth/register", apiCfg.handlerUsersRegister)
	mux.HandleFunc("POST /auth/login", apiCfg.handlerUsersLogin)

	// Create Server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
