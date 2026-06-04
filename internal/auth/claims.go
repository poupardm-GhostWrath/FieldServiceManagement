// internal/auth/claims.go
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, email string, roles []string, secretKey []byte) (string, error) {
	expiresAt := time.Now().UTC().Add(24 * time.Hour)

	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Issuer:    string(TokenTypeAccess),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string, secretKey []byte) (uuid.UUID, error) {
	claimsStruct := CustomClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (any, error) { return secretKey, nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, errors.New("unknown claims type")
	}

	userIDString := claims.UserID

	issuer := claims.Issuer
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %v", err)
	}

	return id, nil
}
