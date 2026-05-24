// internal/handlers/users.go
package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Users []models.User
	}

	// 1. Fetch Users
	dbUsers, err := config.APICfg.DBQueries.GetUsers(r.Context())
	if err != nil {
		http.Error(w, "Couldn't retrieve users", http.StatusInternalServerError)
		return
	}

	users := []models.User{}

	// 2. Get User Roles for each User
	for _, dbUser := range dbUsers {
		dbRoles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Failed to retrieve user roles", http.StatusInternalServerError)
		}
		user := models.User{
			ID:        dbUser.ID.String(),
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Phone:     dbUser.Phone,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Roles:     dbRoles,
		}
		users = append(users, user)
	}

	// 3. Response
	RespondWithJSON(w, http.StatusOK, Response{
		Users: users,
	})
}
