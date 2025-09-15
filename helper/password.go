package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a plain text password with its hashed version
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePasswordStrength validates if password meets minimum requirements
func ValidatePasswordStrength(password string) bool {
	// Minimum 6 characters untuk development
	// Dalam produksi sebaiknya lebih strict
	return len(password) >= 6
}
