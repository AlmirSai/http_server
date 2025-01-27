package validator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	emailRegex        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordMinLength = 8
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < passwordMinLength {
		return &ValidationError{Field: "password", Message: fmt.Sprintf("password must be at least %d characters", passwordMinLength)}
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return &ValidationError{Field: "password", Message: "password must contain at least one uppercase letter, one lowercase letter, and one number"}
	}

	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	if len(name) < 2 {
		return &ValidationError{Field: "name", Message: "name must be at least 2 characters"}
	}
	return nil
}
