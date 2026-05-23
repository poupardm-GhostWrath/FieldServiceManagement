// internal/handlers/auth_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/auth"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
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
	dbUser, err := config.APICfg.DBQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 2. Verify Password
	match, err := argon2id.ComparePasswordAndHash(req.Password, dbUser.PasswordHash)
	if err != nil || !match {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. Get User Roles
	roles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
	if err != nil {
		http.Error(w, "Failed to load roles", http.StatusInternalServerError)
		return
	}

	// 4. Generate Token
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	tokenString, err := auth.GenerateToken(dbUser.ID.String(), dbUser.Email, roles, secretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// 5. Send Response
	RespondWithJSON(w, http.StatusOK, LoginResponse{
		Token: tokenString,
		User: models.User{
			ID:    dbUser.ID.String(),
			Email: dbUser.Email,
			Roles: roles,
		},
	})
}
