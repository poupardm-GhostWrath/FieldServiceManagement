// internal/services/password.go
package services

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	passwordRegex = regexp.MustCompile(`^[A-Za-z0-9@$!%*?&]{12,64}$`)
	specialChars  = "@$!%*?&"
	minLength     = 12
	maxLength     = 64
)

func ValidatePassword(password string) error {
	// 1. Check length
	if len(password) < minLength || len(password) > maxLength {
		return errors.New("password length must be between 12 and 64 characters")
	}

	// 2. Check pattern
	if !passwordRegex.MatchString(password) {
		return errors.New("password contains invalid characters")
	}

	// 3. Check for uppercase
	hasUpper := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return errors.New("password must contain at least one uppercase")
	}

	// 4. Check for lowercase
	hasLower := false
	for _, r := range password {
		if unicode.IsLower(r) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase")
	}

	// 5. Check for digit
	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	// 6. Check for special character
	hasSpecial := false
	for _, r := range password {
		if strings.ContainsRune(specialChars, r) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
