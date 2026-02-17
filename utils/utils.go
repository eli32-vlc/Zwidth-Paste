package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	urlCharset  = "abcdefghijklmnopqrstuvwxyz0123456789"
	codeCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomURL generates a random URL slug
func GenerateRandomURL(length int) string {
	return generateRandomString(urlCharset, length)
}

// GenerateRandomCode generates a random edit/modify code
func GenerateRandomCode(length int) string {
	return generateRandomString(codeCharset, length)
}

// generateRandomString generates a cryptographically random string
func generateRandomString(charset string, length int) string {
	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			// Fallback to a different position if random generation fails
			panic("failed to generate random number: " + err.Error())
		}
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// ValidateURL validates a custom URL
func ValidateURL(url string) bool {
	if len(url) < 2 || len(url) > 100 {
		return false
	}

	// Only allow lowercase letters, numbers, underscores, and hyphens
	matched, _ := regexp.MatchString(`^[a-z0-9_-]+$`, strings.ToLower(url))
	return matched
}

// ValidateCode validates an edit or modify code
func ValidateCode(code string) bool {
	if len(code) < 1 || len(code) > 100 {
		return false
	}
	return true
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSessionID generates a random session ID
func GenerateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetClientIP extracts the client IP from a request (handles proxies)
func GetClientIP(remoteAddr, xForwardedFor, xRealIP string) string {
	if xRealIP != "" {
		return xRealIP
	}
	if xForwardedFor != "" {
		// Take the first IP in the list
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}
	// Extract IP from remoteAddr (format: "IP:port")
	parts := strings.Split(remoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return remoteAddr
}
