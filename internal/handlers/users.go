// internal/handlers/users.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
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
			ID:        dbUser.ID,
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
	pgPhone := services.ValidatePhone(params.Phone)

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
			ID:        dbUser.ID,
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
			ID:        dbUser.ID,
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
	if params.FirstName == "" {
		http.Error(w, "Missing First Name", http.StatusBadRequest)
		return
	}
	if params.LastName == "" {
		http.Error(w, "Missing Last Name", http.StatusBadRequest)
		return
	}
	pgPhone := services.ValidatePhone(params.Phone)

	// 5. Get User From DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't retrieve user", http.StatusInternalServerError)
		return
	}

	// 6. Check if password needs to change
	passwordHash := dbUser.PasswordHash
	if params.Password != "" {
		// Check if password is valid
		if err := services.ValidatePassword(params.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create new password hash
		passwordHash, err = auth.HashPassword(params.Password)
		if err != nil {
			http.Error(w, "Couldn't generate hash", http.StatusInternalServerError)
			return
		}
	}

	// 7. Update User
	dbUserUpdated, err := config.APICfg.DBQueries.UpdateUserProfileByID(r.Context(), database.UpdateUserProfileByIDParams{
		ID:           dbUser.ID,
		PasswordHash: passwordHash,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Phone:        pgPhone,
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
			ID:        dbUserUpdated.ID,
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
			ID:        dbUser.ID,
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

func UserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string   `json:"email"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Phone     string   `json:"phone"`
		Roles     []string `json:"roles"`
	}

	type response struct {
		User models.User `json:"user"`
	}

	// 1. Fetch parameters from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 2. Verify email
	if ok := services.ValidateEmail(params.Email); !ok {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// 3. Generate default password
	password := "DefaultPassword123!"
	hash, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Couldn't generate hash", http.StatusInternalServerError)
		return
	}

	// 4. Check if phone is set
	pgPhone := services.ValidatePhone(params.Phone)

	// 5. Get Roles
	dbRoles := []database.Role{}
	for _, role := range params.Roles {
		dbRole, err := config.APICfg.DBQueries.GetRoleByName(r.Context(), role)
		if err != nil {
			http.Error(w, "Couldn't retrieve role", http.StatusInternalServerError)
			return
		}
		dbRoles = append(dbRoles, dbRole)
	}

	// 6. Create User
	dbUser, err := config.APICfg.DBQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: hash,
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		Phone:        pgPhone,
	})
	if err != nil {
		http.Error(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}

	// 7. Create User Roles
	roles := []string{}
	for _, role := range dbRoles {
		dbUserRole, err := config.APICfg.DBQueries.CreateUserRoles(r.Context(), database.CreateUserRolesParams{
			UserID: dbUser.ID,
			RoleID: role.ID,
		})
		if err != nil {
			http.Error(w, "Couldn't create user role", http.StatusInternalServerError)
			return
		}
		roles = append(roles, dbUserRole.RoleName)
	}

	// 8. Respond
	RespondWithJSON(w, http.StatusCreated, response{
		User: models.User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Phone:     dbUser.Phone,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Roles:     roles,
		},
	})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email     string   `json:"email"`
		FirstName string   `json:"first_name"`
		LastName  string   `json:"last_name"`
		Phone     string   `json:"phone"`
		Roles     []string `json:"roles"`
	}

	type response struct {
		User models.User `json:"user"`
	}

	// 1. Fetch user ID
	userIDString := r.PathValue("userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "Couldn't parse ID", http.StatusBadRequest)
		return
	}

	// 2. Fetch parameters from request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		http.Error(w, "Couldn't decode parameters", http.StatusBadRequest)
		return
	}

	// 3. Get user from DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't get user", http.StatusBadRequest)
		return
	}

	// 4. Verify email
	if ok := services.ValidateEmail(params.Email); !ok {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// 5. Check if phone is set
	pgPhone := services.ValidatePhone(params.Phone)

	// 6. Update User
	dbUser, err = config.APICfg.DBQueries.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		ID:        dbUser.ID,
		Email:     params.Email,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     pgPhone,
	})
	if err != nil {
		http.Error(w, "Couldn't update user", http.StatusInternalServerError)
		return
	}

	// 7. Delete User Roles
	err = config.APICfg.DBQueries.DeleteUserRoles(r.Context(), dbUser.ID)
	if err != nil {
		http.Error(w, "Couldn't delete user roles", http.StatusInternalServerError)
		return
	}

	// 8. Add new user roles
	newRoles := []string{}
	for _, role := range params.Roles {
		dbRole, err := config.APICfg.DBQueries.GetRoleByName(r.Context(), role)
		if err != nil {
			http.Error(w, "Couldn't retrieve roles", http.StatusInternalServerError)
			return
		}
		dbUserRole, err := config.APICfg.DBQueries.CreateUserRoles(r.Context(), database.CreateUserRolesParams{
			UserID: dbUser.ID,
			RoleID: dbRole.ID,
		})
		if err != nil {
			http.Error(w, "Couldn't create user role", http.StatusInternalServerError)
			return
		}
		newRoles = append(newRoles, dbUserRole.RoleName)
	}

	// 9. Respond
	RespondWithJSON(w, http.StatusOK, response{
		User: models.User{
			ID:        dbUser.ID,
			Email:     dbUser.Email,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.LastName,
			Phone:     dbUser.Phone,
			IsActive:  dbUser.IsActive,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Roles:     newRoles,
		},
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch user id
	userIDString := r.PathValue("userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "Couldn't get user ID", http.StatusBadRequest)
		return
	}

	// 2. Get user from DB
	dbUser, err := config.APICfg.DBQueries.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't get user", http.StatusInternalServerError)
		return
	}

	// 3. Check if user already "deleted"
	if !dbUser.IsActive.Bool {
		http.Error(w, "User already deleted", http.StatusBadRequest)
		return
	}

	// 4. Delete User From DB (Soft Delete)
	err = config.APICfg.DBQueries.DeleteUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't delete user", http.StatusInternalServerError)
		return
	}

	// 5. Delete User Roles
	err = config.APICfg.DBQueries.DeleteUserRoles(r.Context(), userID)
	if err != nil {
		http.Error(w, "Couldn't delete user roles", http.StatusInternalServerError)
		return
	}

	// 6. Respond
	w.WriteHeader(http.StatusNoContent)
}
