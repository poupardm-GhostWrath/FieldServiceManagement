// internal/services/email_test.go
package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	// Test: Good simple email
	ok := ValidateEmail("admin@example.com")
	assert.True(t, ok)

	// Test: Good '.' compound email
	ok = ValidateEmail("admin.test@example.com")
	assert.True(t, ok)

	// Test: Good '+tag' compound email
	ok = ValidateEmail("admin+test@example.com")
	assert.True(t, ok)

	// Test: Good compound domain
	ok = ValidateEmail("admin@test.example.com")
	assert.True(t, ok)

	// Test: Good multiple Top-Level Domain
	ok = ValidateEmail("admin@example.gov.com")
	assert.True(t, ok)

	// Test: Bad trailing '.'
	ok = ValidateEmail(".invalid@example.com")
	assert.False(t, ok)

	// Test: Bad Missing domain
	ok = ValidateEmail("admin@.com")
	assert.False(t, ok)

	// Test: Bad Missing TLD
	ok = ValidateEmail("admin@example")
	assert.False(t, ok)

	// Test: Bad Missing '@'
	ok = ValidateEmail("adminexample.com")
	assert.False(t, ok)
}
