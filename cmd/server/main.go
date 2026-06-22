// Package main
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/handlers"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/middleware"
)

func main() {
	// Get Environmental Data
	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if filepathRoot == "" {
		log.Fatal("FILEPATH_ROOT must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	// Create Router
	r := chi.NewRouter()

	// Set Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Public Routes
	r.Route("/auth", func(r chi.Router) {
		r.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/login", handlers.Login)          // Login User (5 req/min limit)
		r.With(httprate.LimitByIP(3, 1*time.Hour)).Post("/register", handlers.RegisterUsers) // Register User (3 req/hour limit)
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware())            // Authentication
		r.Use(httprate.LimitByIP(120, 1*time.Minute)) // 120 req/min limit

		// User Management
		r.Get("/users/me", handlers.GetUserProfile)                                                           // Get Current User Profile
		r.Put("/users/me", handlers.UpdateUserProfile)                                                        // Update Current User Profile
		r.With(middleware.RequireRole("admin", "dispatcher")).Get("/users", handlers.ListUsers)               // List all users
		r.With(middleware.RequireRole("admin", "dispatcher")).Get("/users/{userID}", handlers.GetUserDetails) // Get specific user details
		r.With(middleware.RequireRole("admin")).Post("/users", handlers.UserCreate)                           // Create New User (Admin)
		r.With(middleware.RequireRole("admin")).Put("/users/{userID}", handlers.UpdateUser)                   // Update User details/roles
		r.With(middleware.RequireRole("admin")).Delete("/users/{userID}", handlers.DeleteUser)                // Delete User (Soft Delete)

		// Customer Management
		r.With(middleware.RequireRole("admin", "dispatcher", "technician")).Get("/customers", handlers.GetCustomers)             // List Customers
		r.With(middleware.RequireRole("admin", "dispatcher", "technician")).Get("/customers/{customerID}", handlers.GetCustomer) // Get Customer details
		r.With(middleware.RequireRole("admin", "dispatcher")).Post("/customers", handlers.CreateCustomer)                        // Create Customer
		r.With(middleware.RequireRole("admin")).Delete("/customers/{customerID}", handlers.DeleteCustomer)                       // Delete customer (soft delete)
		r.With(middleware.RequireRole("admin", "dispatcher")).Put("/customers/{customerID}", handlers.UpdateCustomer)            // Update Customer

		// Inventory Management
		r.With(middleware.RequireRole("admin", "dispatcher")).Post("/inventory", handlers.CreateInventoryItem)                 // Create Inventory Item
		r.With(middleware.RequireRole("admin", "dispatcher", "technician")).Get("/inventory", handlers.GetInventory)           // Get Inventory
		r.With(middleware.RequireRole("admin", "dispatcher", "technician")).Get("/inventory/{sku}", handlers.GetInventoryItem) // Get Inventory Item
		r.With(middleware.RequireRole("admin", "dispatcher")).Put("/inventory/{sku}", handlers.UpdateInventoryItem)            // Update Inventory Item
	})

	// Initial Entrypoint
	fs := http.FileServer(http.Dir(filepathRoot))
	r.Mount("/", http.StripPrefix("/", fs))

	defer config.APICfg.DB.Close(context.Background())
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
