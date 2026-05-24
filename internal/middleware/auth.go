// internal/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/auth"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/config"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			secretKey := []byte(os.Getenv("JWT_SECRET"))

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(token *jwt.Token) (any, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secretKey, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			dbUser, err := config.APICfg.DBQueries.GetUserByEmail(r.Context(), claims.Email)
			if err != nil {
				http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
				return
			}

			dbUserRole, err := config.APICfg.DBQueries.GetUserRoles(r.Context(), dbUser.ID)
			if err != nil {
				http.Error(w, "Failed to retrieve user roles", http.StatusInternalServerError)
				return
			}

			user := models.User{
				ID:    dbUser.ID.String(),
				Email: dbUser.Email,
				Roles: dbUserRole,
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve user from context (set by AuthMiddleware)
			user, ok := r.Context().Value(userContextKey).(*models.User)
			if !ok || user == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, userRole := range user.Roles {
				if slices.Contains(requiredRoles, userRole) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
