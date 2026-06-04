// Testing main
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/auth"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/database"
)

func main() {
	hash, err := auth.HashPassword("admin")
	if err != nil {
		log.Fatalf("Failed to hash password: %v\n", err)
	}

	fmt.Println("Creating test admin user")

	user, err := config.APICfg.DBQueries.CreateUser(context.Background(), database.CreateUserParams{
		Email:        "admin@example.com",
		PasswordHash: hash,
		FirstName:    "admin_f",
		LastName:     "admin_l",
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
	role, err := config.APICfg.DBQueries.GetRoleByName(context.Background(), "admin")
	if err != nil {
		log.Fatalf("failed to get admin role: %v\n", err)
	}
	fmt.Println("===== Role =====")
	fmt.Printf("ID: %v\n", role.ID)
	fmt.Printf("Name: %s\n", role.Name)

	fmt.Println()
	fmt.Println("Creating user role entry")
	userRole, err := config.APICfg.DBQueries.CreateUserRoles(context.Background(), database.CreateUserRolesParams{
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
