// internal/handlers/users.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/auth"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Users []models.User `json:"users"`
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
	RespondWithJSON(w, http.StatusOK, response{
		Users: users,
	})
}

func RegisterUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	type response struct {
		User models.User `json:"user"`
	}

	// 1. Fetch data from Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	// 2. Create Password Hash
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, "Couldn't generate hash", http.StatusInternalServerError)
		return
	}

	// 3. Check if phone is set
	pgPhone := pgtype.Text{
		String: params.Phone,
		Valid:  true,
	}
	if params.Phone == "" {
		pgPhone.Valid = false
	}

	// 4. Create DB User
	dbUser, err := config.APICfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: passwordHash,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Phone:        pgPhone,
	})
	if err != nil {
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}

	// 5. Create Default User Role
	dbRole, err := config.APICfg.DBQueries.CreateUserRoles(r.Context(), database.CreateUserRolesParams{
		UserID: dbUser.ID,
		RoleID: 4,
	})
	if err != nil {
		http.Error(w, "Couldn't create user role", http.StatusInternalServerError)
		return
	}

	// 6. Respond
	RespondWithJSON(w, http.StatusCreated, response{
		User: models.User{
			ID:        dbUser.ID.String(),
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Phone:     dbUser.Phone,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Roles:     []string{dbRole.RoleName},
		},
	})
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User models.User `json:"user"`
	}

	// 1. Get Token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Verify Token
	userID, err := auth.ValidateToken(token, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Get User from DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	// 4. Get User Roles
	dbRoles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user roles", http.StatusInternalServerError)
		return
	}

	// 5. Respond
	RespondWithJSON(w, http.StatusOK, response{
		User: models.User{
			ID:        dbUser.ID.String(),
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Phone:     dbUser.Phone,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Roles:     dbRoles,
		},
	})
}
