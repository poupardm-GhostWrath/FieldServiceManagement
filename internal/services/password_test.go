// internal/services/password_test.go
package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePassword(t *testing.T) {
	// Test: Good password
	err := ValidatePassword("StrongP@ssw0rd")
	require.NoError(t, err)

	// Test: Short Password
	err = ValidatePassword("Sh0rt!")
	require.Error(t, err)

	// Test: Long Password
	err = ValidatePassword("ThisIsTheLongestPasswordEver!ThisIsTheLongestPasswordEver!ThisIsTheLongestPasswordEver!ThisIsTheLongestPasswordEver!ThisIsTheLongestPasswordEver!ThisIsTheLongestPasswordEver!")
	require.Error(t, err)

	// Test: Invalid character
	err = ValidatePassword("Passw0rd Something!")
	require.Error(t, err)

	// Test: Missing uppercase
	err = ValidatePassword("nouppercase1@")
	require.Error(t, err)

	// Test: Missing lowercase
	err = ValidatePassword("NOLOWERCASE1@")
	require.Error(t, err)

	// Test: Missing digit
	err = ValidatePassword("P@sswordNoDigit")
	require.Error(t, err)

	// Test: Missing special character
	err = ValidatePassword("NoSpec1alCharacter")
	require.Error(t, err)
}
