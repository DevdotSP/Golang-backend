package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// generateVerificationToken creates a random verification token
func GenerateVerificationToken() string {
	b := make([]byte, 16) // 16 bytes for a 32-character hex string
	rand.Read(b)
	return hex.EncodeToString(b)
}
