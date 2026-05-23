// internal/handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Fetch user from DB
}
