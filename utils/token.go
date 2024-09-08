package utils

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define custom error messages
var (
	ErrInvalidSigningMethod = errors.New("invalid signing method")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenNotFound        = errors.New("token not found")
)

// ActiveTokens to manage active tokens
var activeTokens = sync.Map{}

// ValidateToken validates the token and returns claims if valid
func ValidateToken(token string) (jwt.MapClaims, error) {
	// Check if the token is blacklisted
	if _, ok := activeTokens.Load(token); !ok {
		return nil, ErrTokenNotFound
	}

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}
		return []byte("your_secret_key"), nil // Replace with your secret key
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateJWT creates a new JWT token, removing the current token if provided
func GenerateJWT(userID uint, currentToken string) (string, error) {
	// If a current token is provided, remove it from the active tokens map
	if currentToken != "" {
		activeTokens.Delete(currentToken)
	}

	// Define the token claims, including a unique claim
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":     time.Now().Unix(),                     // Issued at
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Replace "your_secret_key" with a strong secret key
	secretKey := []byte("your_secret_key")

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Store the new token in the active tokens map
	activeTokens.Store(signedToken, struct{}{})

	return signedToken, nil
}

// DeleteToken removes a token from the active tokens map
func DeleteToken(token string) error {
	_, loaded := activeTokens.LoadAndDelete(token)
	if !loaded {
		return ErrTokenNotFound
	}
	log.Printf("Token %s has been removed from active tokens.", token)
	return nil
}
