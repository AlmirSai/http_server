package validator

import (
	"fmt"
	"unicode"
)

type PasswordConfig struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
}

type PasswordValidator struct {
	config PasswordConfig
}

func NewPasswordValidator(config PasswordConfig) *PasswordValidator {
	return &PasswordValidator{config: config}
}

func (v *PasswordValidator) Validate(password string) error {
	if len(password) < v.config.MinLength {
		return fmt.Errorf("password must be at least %d characters long", v.config.MinLength)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if v.config.RequireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if v.config.RequireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if v.config.RequireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if v.config.RequireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}
