// Reset databases
package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// Check if in dev environment
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	err := cfg.reset()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't reset database", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset to initial state"))
}

func (cfg *apiConfig) reset() error {
	// Purge audit_logs
	if _, err := cfg.db.Exec("DELETE FROM audit_logs"); err != nil {
		return fmt.Errorf("failed to reset table audit_logs: %v", err)
	}
	// Purge invoices
	if _, err := cfg.db.Exec("DELETE FROM invoices"); err != nil {
		return fmt.Errorf("failed to reset table invoices: %v", err)
	}
	// Purge job_labor
	if _, err := cfg.db.Exec("DELETE FROM job_labor"); err != nil {
		return fmt.Errorf("failed to reset table job_labor: %v", err)
	}
	// Purge job_parts
	if _, err := cfg.db.Exec("DELETE FROM job_parts"); err != nil {
		return fmt.Errorf("failed to reset table job_parts: %v", err)
	}
	// Purge jobs
	if _, err := cfg.db.Exec("DELETE FROM jobs"); err != nil {
		return fmt.Errorf("failed to reset table jobs: %v", err)
	}
	// Purge inventory_items
	if _, err := cfg.db.Exec("DELETE FROM inventory_items"); err != nil {
		return fmt.Errorf("failed to reset table inventory_items", err)
	}
	// Purge customers
	if _, err := cfg.db.Exec("DELETE FROM customers"); err != nil {
		return fmt.Errorf("failed to reset table customers", err)
	}
	// Purge user_roles
	if _, err := cfg.db.Exec("DELETE FROM user_roles"); err != nil {
		return fmt.Errorf("failed to reset table user_roles", err)
	}
	// Purge users
	if _, err := cfg.db.Exec("DELETE FROM users"); err != nil {
		return fmt.Errorf("failed to reset table users", err)
	}
	return nil
}
