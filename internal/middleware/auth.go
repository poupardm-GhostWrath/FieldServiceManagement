// internal/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/models"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			user := &models.User{
				ID:    claims["user_id"].(string),
				Roles: parseRolesFromClaims(claims),
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

func parseRolesFromClaims(claims jwt.MapClaims) []string {
	var roles []string
	if roleData, ok := claims["roles"].([]any); ok {
		for _, r := range roleData {
			if roleName, ok := r.(string); ok {
				roles = append(roles, roleName)
			}
		}
	}
	return roles
}
