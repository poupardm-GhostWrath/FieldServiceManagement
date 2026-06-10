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

	// Initial Entrypoint
	r.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// Public Routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handlers.Login)            // Login User
		r.Post("/register", handlers.RegisterUsers) // Register User
	})

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware()) // Authentication

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
	})

	// Customers
	/*
		mux.HandleFunc("GET /customers", apiCfg.handlerCustomersGet)												// Admin, Dispatcher, Technician
		mux.HandleFunc("GET /customers/{customerID}", apiCfg.handlerCustomersGetByID)				// Admin, Dispatcher, Technician, Customer
		mux.HandleFunc("POST /customers", apiCfg.handlerCustomersCreate)										// Admin, Dispatcher
		mux.HandleFunc("PUT /customers/{customerID}", apiCfg.handlerCustomersUpdate)				// Admin, Dispatcher
		mux.HandleFunc("DELETE /customers/{customerID}", apiCfg.handlerCustomersDelete)			// Admin (Soft delete)
		mux.HandleFunc("GET /customers/{customerID}/jobs", apiCfg.handlerCustomersGetJobs)	// Admin, Dispatcher, Technician, Customer
	*/

	// Inventory
	/*
		mux.HandleFunc("GET /inventory", apiCfg.handlerIventoryGet)													// Admin, Dispatcher, Technician
		mux.HandleFunc("GET /inventory/{itemID}", apiCfg.handlerInventoryGetByID)						// Admin, Dispatcher, Technician
		mux.HandleFunc("POST /inventory", apiCfg.handlerInventoryCreate)										// Admin, Dispatcher
		mux.HandleFunc("PUT /inventory/{itemID}", apiCfg.handlerInventoryUpdate)						// Admin, Dispatcher
		mux.HandleFunc("DELETE /inventory/{itemID}", apiCfg.handlerInventoryDelete)					// Admin
		mux.HandleFunc("POST /inventory/{itemID}/restock", apiCfg.handlerInventoryRestock)	// Admin, Dispatcher
		mux.HandleFunc("GET /inventory/alerts", apiCfg.handlerInventoryAlert)								// Admin, Dispatcher
	*/

	// Jobs
	/*
		mux.HandleFunc("GET /jobs", apiCfg.handlerJobsGet)																		// Admin, Dispatcher, Technician
		mux.HandleFunc("GET /jobs/{jobID}", apiCfg.handlerJobsGetByID)												// Admin, Dispatcher, Technician
		mux.HandleFunc("POST /jobs", apiCfg.handlerJobsCreate)																// Admin, Dispatcher
		mux.HandleFunc("PUT /jobs/{jobID}", apiCfg.handlerJobsUpdate)													// Admin, Dispatcher
		mux.HandleFunc("DELETE /jobs/{jobID}", apiCfg.handlerJobsDelete)											// Admin, Dispatcher
		mux.HandleFunc("PATCH /jobs/{jobID}/status", apiCfg.handlerJobsPatch)									// Admin, Dispatcher, Technician
		mux.HandleFunc("POST /jobs/{jobID}/parts", apiCfg.handlerJobsAddParts)								// Dispatcher, Technician
		mux.HandleFunc("POST /jobs/{jobID}/labor", apiCfg.handlerJobsLogLabor)								// Dispatcher, Technician
		mux.HandleFunc("GET /jobs/schedule", apiCfg.handlerJobsGetSchedule)										// Admin, Dispatcher
		mux.HandleFunc("GET /jobs/tech/{techID}/schedule", apiCfg.handlerJobsGetTechSchedule)	// Admin, Dispatcher, Technician
	*/

	// Invoices
	/*
		mux.HandleFunc("GET /invoices", apiCfg.handlerInvoicesGet)										// Admin, Dispatcher, Technician, Customer
		mux.HandleFunc("GET /invoices/{invoiceID}", apiCfg.handlerInvoicesGetByID)		// Admin, Dispatcher, Technician, Customer
		mux.HandleFunc("POST /invoices", apiCfg.handlerInvoicesCreate)								// Admin, Dispatcher
		mux.HandleFunc("PUT /invoices/{invoiceID}", apiCfg.handlerInvoicesUpdate)			// Admin
		mux.HandleFunc("POST /invoices/{invoiceID}/pay", apiCfg.handlerInvoicesPaid)	// Admin, Customer
		mux.HandleFunc("GET /invoices/{invoiceID}/pdf", apiCfg.handlerInvoicesPDF)		// Admin, Dispatcher, Technician, Customer
	*/

	// Reports & Analytics
	/*
		mux.HandleFunc("GET /reports/revenue", apiCfg.handlerReportsRevenue)									// Admin
		mux.HandleFunc("GET /reports/technician-performance", apiCfg.handlerReportsTechPerf)	// Admin, Dispatcher
		mux.HandleFunc("GET /reports/inventory-usage", apiCfg.handlerReportsInventory)				// Admin
	*/
	defer config.APICfg.DB.Close(context.Background())
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
