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
		RespondWithError(w, http.StatusBadRequest, "Couldn't decode request", err)
		return
	}

	// 1. Fetch user from DB
	dbUser, err := config.APICfg.DBQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "User not found", err)
		return
	}

	// 2. Verify Password
	match, err := argon2id.ComparePasswordAndHash(req.Password, dbUser.PasswordHash)
	if err != nil || !match {
		RespondWithError(w, http.StatusBadRequest, "Invalid credentials", err)
		return
	}

	// 3. Get User Roles
	roles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user roles", err)
		return
	}

	// 4. Generate Token
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		RespondWithError(w, http.StatusInternalServerError, "Server configuration error", nil)
		return
	}

	tokenString, err := auth.GenerateToken(dbUser.ID.String(), dbUser.Email, roles, secretKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	// 5. Send Response
	RespondWithJSON(w, http.StatusOK, LoginResponse{
		Token: tokenString,
		User: models.User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Roles:     roles,
		},
	})
}
