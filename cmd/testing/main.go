// Testing main
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/middleware"
)

func main() {
	const dbURL = "postgres://postgres:postgres@localhost:5432/fieldservicemanagement?sslmode=disable"

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close(context.Background())

	hash, err := middleware.HashPassword("admin")
	if err != nil {
		log.Fatalf("Failed to hash password: %v\n", err)
	}

	dbQueries := database.New(db)

	fmt.Println("Creating test admin user")

	user, err := dbQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:        "test@example.com",
		PasswordHash: hash,
		FirstName:    "test",
		LastName:     "admin",
	})
	if err != nil {
		log.Fatalf("failed to create user: %v\n", err)
	}

	fmt.Println("===== User =====")
	fmt.Printf("ID: %s\n", user.ID.String())
	fmt.Printf("First Name: %s\n", user.FirstName)
	fmt.Printf("Last Name: %s\n", user.LastName)

	fmt.Println()
	fmt.Println("Getting admin role")
	role, err := dbQueries.GetRoleByName(context.Background(), "admin")
	if err != nil {
		log.Fatalf("failed to get admin role: %v\n", err)
	}
	fmt.Println("===== Role =====")
	fmt.Printf("ID: %v\n", role.ID)
	fmt.Printf("Name: %s\n", role.Name)

	fmt.Println()
	fmt.Println("Creating user role entry")
	userRole, err := dbQueries.CreateUserRoles(context.Background(), database.CreateUserRolesParams{
		UserID: user.ID,
		RoleID: role.ID,
	})
	if err != nil {
		log.Fatalf("failed to create user role: %v\n", err)
	}
	fmt.Println("===== User Role =====")
	fmt.Printf("User ID: %s\n", userRole.UserID.String())
	fmt.Printf("Role ID: %v\n", userRole.RoleID)
}
