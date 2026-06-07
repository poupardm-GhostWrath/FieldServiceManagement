// internal/handlers/users.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/auth"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/services"
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

	// 2. Verify email
	valid := services.ValidateEmail(params.Email)
	if !valid {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// 3. Verify Password
	err = services.ValidatePassword(params.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Create Password Hash
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, "Couldn't generate hash", http.StatusInternalServerError)
		return
	}

	// 4. Check if phone is set
	pgPhone := pgtype.Text{
		String: params.Phone,
		Valid:  true,
	}
	if params.Phone == "" {
		pgPhone.Valid = false
	}

	// 5. Create DB User
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

	// 6. Create Default User Role
	dbRole, err := config.APICfg.DBQueries.CreateUserRoles(r.Context(), database.CreateUserRolesParams{
		UserID: dbUser.ID,
		RoleID: 4,
	})
	if err != nil {
		http.Error(w, "Couldn't create user role", http.StatusInternalServerError)
		return
	}

	// 7. Respond
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

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
	}

	type response struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
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

	// 3. Decode Request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 4. Check for invalid data
	if valid := services.ValidateEmail(params.Email); !valid {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if err := services.ValidatePassword(params.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if params.FirstName == "" {
		http.Error(w, "Missing First Name", http.StatusBadRequest)
		return
	}
	if params.LastName == "" {
		http.Error(w, "Missing Last Name", http.StatusBadRequest)
		return
	}

	// 5. Create new password hash
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, "Couldn't generate hash", http.StatusInternalServerError)
		return
	}

	// 6. Get User From DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	// 7. Update User
	dbUserUpdated, err := config.APICfg.DBQueries.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:           dbUser.ID,
		Email:        params.Email,
		PasswordHash: passwordHash,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Phone: pgtype.Text{
			String: params.Phone,
			Valid:  (params.Phone != ""),
		},
	})
	if err != nil {
		http.Error(w, "Couldn't update user", http.StatusInternalServerError)
		return
	}

	// 8. Get Roles
	roles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUserUpdated.ID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user roles", http.StatusInternalServerError)
		return
	}

	// 9. Generate New Token
	token, err = auth.GenerateToken(dbUserUpdated.ID.String(), dbUserUpdated.Email, roles, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Couldn't generate token", http.StatusInternalServerError)
		return
	}

	// 9. Respond
	RespondWithJSON(w, http.StatusOK, response{
		Token: token,
		User: models.User{
			ID:        dbUserUpdated.ID.String(),
			Email:     dbUserUpdated.Email,
			FirstName: dbUserUpdated.FirstName,
			LastName:  dbUserUpdated.LastName,
			Phone:     dbUserUpdated.Phone,
			IsActive:  dbUserUpdated.IsActive,
			CreatedAt: dbUserUpdated.CreatedAt,
			UpdatedAt: dbUserUpdated.UpdatedAt,
			Roles:     roles,
		},
	})
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User models.User `json:"user"`
	}

	// 1. Get request user id
	userIDString := r.PathValue("userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 2. Get user from DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	// 3. Get user roles from DB
	dbRoles, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user roles", http.StatusInternalServerError)
		return
	}

	// 3. Respond
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
