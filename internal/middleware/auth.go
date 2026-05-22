// internal/middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/poupardm-GhostWrath/FieldServiceManagement/internal/handlers"
)

type (
	TokenType  string
	contextKey string
)

const (
	TokenTypeAccess  TokenType  = "field-service-management-access"
	userIDContextKey contextKey = "userID"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				handlers.RespondWithError(w, http.StatusUnauthorized, "Missing authorization header", nil)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate token
			claimsStruct := jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(
				tokenString,
				&claimsStruct,
				func(token *jwt.Token) (any, error) { return secretKey, nil },
			)
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "failed to parse token", err)
				return
			}

			userIDString, err := token.Claims.GetSubject()
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "failed to get user ID", err)
				return
			}

			issuer, err := token.Claims.GetIssuer()
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "failed to get issuer", err)
				return
			}
			if issuer != string(TokenTypeAccess) {
				handlers.RespondWithError(w, http.StatusUnauthorized, "invalid issuer", nil)
				return
			}

			id, err := uuid.Parse(userIDString)
			if err != nil {
				handlers.RespondWithError(w, http.StatusUnauthorized, "invalid user ID", err)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
