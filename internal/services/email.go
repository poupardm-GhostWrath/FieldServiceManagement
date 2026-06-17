// internal/services/email.go
package services

import (
	"regexp"
)

var emailRegex = regexp.MustCompile(`[A-Za-z0-9]([A-Za-z0-9._%+-]*[A-Za-z0-9])?@[A-Za-z0-9]([A-Za-z0-9.-]*[A-Za-z0-9])?\.[A-Za-z]{2,}`)

func ValidateEmail(email string) bool {
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}
