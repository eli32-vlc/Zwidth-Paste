package hashcash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBits is the default number of leading zero bits required
	DefaultBits = 20
	// ValidityDuration is how long a hashcash is valid for
	ValidityDuration = 5 * time.Minute
)

// Challenge represents a hashcash challenge
type Challenge struct {
	Resource  string
	Bits      int
	Timestamp int64
	Rand      string
	Counter   int
}

// NewChallenge creates a new hashcash challenge
func NewChallenge(resource string, bits int) *Challenge {
	return &Challenge{
		Resource:  resource,
		Bits:      bits,
		Timestamp: time.Now().Unix(),
		Rand:      generateRandomString(16),
		Counter:   0,
	}
}

// String returns the hashcash string representation
func (c *Challenge) String() string {
	return fmt.Sprintf("1:%d:%d:%s::%s:%x",
		c.Bits, c.Timestamp, c.Resource, c.Rand, c.Counter)
}

// Parse parses a hashcash string
func Parse(hashcash string) (*Challenge, error) {
	parts := strings.Split(hashcash, ":")
	if len(parts) != 7 {
		return nil, fmt.Errorf("invalid hashcash format")
	}

	if parts[0] != "1" {
		return nil, fmt.Errorf("unsupported hashcash version")
	}

	bits, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid bits: %w", err)
	}

	timestamp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	counter, err := strconv.ParseInt(parts[6], 16, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid counter: %w", err)
	}

	return &Challenge{
		Resource:  parts[3],
		Bits:      bits,
		Timestamp: timestamp,
		Rand:      parts[5],
		Counter:   int(counter),
	}, nil
}

// Verify verifies a hashcash solution
func Verify(hashcash string, resource string, requiredBits int) error {
	challenge, err := Parse(hashcash)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// Check if resource matches
	if challenge.Resource != resource {
		return fmt.Errorf("resource mismatch")
	}

	// Check if bits requirement is met
	if challenge.Bits < requiredBits {
		return fmt.Errorf("insufficient bits: got %d, required %d", challenge.Bits, requiredBits)
	}

	// Check if not expired
	age := time.Now().Unix() - challenge.Timestamp
	if age > int64(ValidityDuration.Seconds()) {
		return fmt.Errorf("hashcash expired")
	}

	// Check if timestamp is not in the future
	if challenge.Timestamp > time.Now().Unix() {
		return fmt.Errorf("timestamp in future")
	}

	// Verify the hash has required leading zero bits
	hash := sha256.Sum256([]byte(hashcash))
	hashHex := hex.EncodeToString(hash[:])

	leadingZeros := countLeadingZeroBits(hashHex)
	if leadingZeros < challenge.Bits {
		return fmt.Errorf("invalid proof of work: got %d bits, required %d", leadingZeros, challenge.Bits)
	}

	return nil
}

// countLeadingZeroBits counts the number of leading zero bits in a hex string
func countLeadingZeroBits(hexStr string) int {
	count := 0
	for _, char := range hexStr {
		val, _ := strconv.ParseInt(string(char), 16, 64)
		if val == 0 {
			count += 4
		} else {
			// Count leading zeros in this nibble
			for i := 3; i >= 0; i-- {
				if (val & (1 << i)) == 0 {
					count++
				} else {
					return count
				}
			}
			return count
		}
	}
	return count
}

// generateRandomString generates a random string of given length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateChallenge generates a new challenge for the client
func GenerateChallenge(resource string) string {
	challenge := NewChallenge(resource, DefaultBits)
	return challenge.String()
}
