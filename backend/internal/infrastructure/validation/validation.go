package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidatePositiveInt validates that an integer is positive (> 0)
func ValidatePositiveInt(val int, fieldName string) error {
	if val <= 0 {
		return fmt.Errorf("%s must be a positive integer, got %d", fieldName, val)
	}
	return nil
}

// ValidateNonNegativeInt validates that an integer is non-negative (>= 0)
func ValidateNonNegativeInt(val int, fieldName string) error {
	if val < 0 {
		return fmt.Errorf("%s must be non-negative, got %d", fieldName, val)
	}
	return nil
}

// ValidateRequired validates that a string is not empty
func ValidateRequired(val string, fieldName string) error {
	if strings.TrimSpace(val) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateEmail validates that an email address is properly formatted
func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email regex pattern
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateMaxLength validates that a string does not exceed max length
func ValidateMaxLength(val string, fieldName string, maxLen int) error {
	if len(val) > maxLen {
		return fmt.Errorf("%s must not exceed %d characters, got %d", fieldName, maxLen, len(val))
	}
	return nil
}
