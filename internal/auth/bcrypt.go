package auth

import (
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
)

var (
	ErrPasswordTooShort      = errors.New("le mot de passe doit contenir au moins 8 caractères")
	ErrPasswordNoUpper       = errors.New("le mot de passe doit contenir au moins une lettre majuscule")
	ErrPasswordNoLower       = errors.New("le mot de passe doit contenir au moins une lettre minuscule")
	ErrPasswordNoDigit       = errors.New("le mot de passe doit contenir au moins un chiffre")
	ErrPasswordNoSpecialChar = errors.New("le mot de passe doit contenir au moins un caractère spécial (@$!%*?&)")
)

// ValidatePassword validates a password according to security requirements
// Requirements:
// - Minimum 8 characters
// - At least one uppercase letter
// - At least one lowercase letter
// - At least one digit
// - At least one special character (@$!%*?&)
func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	specialChars := "@$!%*?&"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasDigit {
		return ErrPasswordNoDigit
	}
	if !hasSpecial {
		return ErrPasswordNoSpecialChar
	}

	return nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

