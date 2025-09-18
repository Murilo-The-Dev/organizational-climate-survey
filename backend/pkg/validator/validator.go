// Package validator provides validation functions for common business entities
// and input data with consistent error handling and comprehensive validation rules.
package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Constants define validation constraints and limits
const (
	MinPasswordLength = 8  // Minimum password length for security compliance
	CNPJLength       = 14 // CNPJ must contain exactly 14 digits
)

// Pre-compiled regex patterns for performance optimization
var (
	// emailRegex validates standard email format according to RFC 5322
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// cnpjRegex removes all non-digit characters for normalization
	cnpjRegex  = regexp.MustCompile(`\D`)
)

// ValidationError represents a structured validation error with field context
type ValidationError struct {
	Field   string // Name of the field that failed validation
	Message string // Human-readable error message
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// IsEmail validates if the provided string is a well-formed email address.
// It trims whitespace and checks against RFC 5322 format.
// Returns ValidationError if email is empty or malformed.
func IsEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return ValidationError{"email", "required"}
	}
	
	if !emailRegex.MatchString(email) {
		return ValidationError{"email", "invalid format"}
	}
	return nil
}

// IsCNPJ validates Brazilian CNPJ (Cadastro Nacional da Pessoa Jurídica) format.
// Performs complete validation including check digit algorithm.
// Returns ValidationError if CNPJ format or check digits are invalid.
func IsCNPJ(cnpj string) error {
	cnpjNumbers := cnpjRegex.ReplaceAllString(cnpj, "")
	
	if len(cnpjNumbers) != CNPJLength {
		return ValidationError{"cnpj", "must contain 14 digits"}
	}
	
	if isAllSameDigits(cnpjNumbers) {
		return ValidationError{"cnpj", "invalid sequence"}
	}
	
	if !isValidCNPJCheckDigits(cnpjNumbers) {
		return ValidationError{"cnpj", "invalid check digits"}
	}
	
	return nil
}

// isAllSameDigits checks if all characters in the string are identical.
// Used to reject invalid CNPJ patterns like "11111111111111".
func isAllSameDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for _, char := range s {
		if byte(char) != first {
			return false
		}
	}
	return true
}

// isValidCNPJCheckDigits implements the CNPJ check digit algorithm.
// Uses official Brazilian algorithm with two verification digits.
func isValidCNPJCheckDigits(cnpj string) bool {
	// Weight arrays for check digit calculation
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	
	// Convert string digits to integer array
	digits := make([]int, CNPJLength)
	for i, char := range cnpj {
		digits[i], _ = strconv.Atoi(string(char))
	}
	
	// Calculate first check digit
	sum := 0
	for i := 0; i < 12; i++ {
		sum += digits[i] * weights1[i]
	}
	remainder := sum % 11
	checkDigit1 := 0
	if remainder >= 2 {
		checkDigit1 = 11 - remainder
	}
	
	if digits[12] != checkDigit1 {
		return false
	}
	
	// Calculate second check digit
	sum = 0
	for i := 0; i < 13; i++ {
		sum += digits[i] * weights2[i]
	}
	remainder = sum % 11
	checkDigit2 := 0
	if remainder >= 2 {
		checkDigit2 = 11 - remainder
	}
	
	return digits[13] == checkDigit2
}

// IsPasswordStrong validates password strength according to security best practices.
// Enforces minimum length and requires uppercase, lowercase, numeric, and special characters.
// Returns ValidationError with specific requirement that failed.
func IsPasswordStrong(password string) error {
	if len(password) < MinPasswordLength {
		return ValidationError{"password", fmt.Sprintf("must be at least %d characters", MinPasswordLength)}
	}
	
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	
	// Check each character type requirement
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	
	// Return specific error for first missing requirement
	if !hasUpper {
		return ValidationError{"password", "must contain uppercase letter"}
	}
	if !hasLower {
		return ValidationError{"password", "must contain lowercase letter"}
	}
	if !hasNumber {
		return ValidationError{"password", "must contain number"}
	}
	if !hasSpecial {
		return ValidationError{"password", "must contain special character"}
	}
	
	return nil
}

// IsValidStatus checks if the provided status is in the list of valid statuses.
// Returns ValidationError with available options if status is invalid.
func IsValidStatus(status string, validStatuses []string) error {
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return ValidationError{"status", fmt.Sprintf("must be one of: %v", validStatuses)}
}