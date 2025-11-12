package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateClientID generates a unique 8-digit numeric client ID
func GenerateClientID() (string, error) {
	// Generate a random 8-digit number
	max := big.NewInt(99999999) // 8 digits max
	min := big.NewInt(10000000) // 8 digits min (starts with 1)

	// Generate random number between min and max
	n, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
	if err != nil {
		return "", err
	}

	// Add min to get the final number
	clientID := new(big.Int).Add(n, min)

	return fmt.Sprintf("%08d", clientID.Int64()), nil
}

// ValidateClientID validates that a client ID is a valid 8-digit number
func ValidateClientID(clientID string) bool {
	if len(clientID) != 8 {
		return false
	}

	// Check if all characters are digits
	for _, char := range clientID {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}
















