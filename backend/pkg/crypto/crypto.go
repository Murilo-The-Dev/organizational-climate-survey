// Package crypto provides secure cryptographic operations including password hashing,
// token generation, and secure comparisons using industry-standard algorithms.
package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Constants define security parameters and constraints
const (
	DefaultBcryptCost = 12 // Recommended bcrypt cost for 2024+ (balance of security/performance)
	MinTokenBytes     = 32 // Minimum token length for cryptographic security
)

// Hasher interface defines password hashing operations
// Abstracts the underlying hashing mechanism for testability and flexibility
type Hasher interface {
	HashPassword(password string) (string, error) // Generate secure hash from password
	CheckPasswordHash(password, hash string) bool // Verify password against hash
}

// TokenGenerator interface defines secure token and random byte generation
// Used for session tokens, API keys, and cryptographic nonces
type TokenGenerator interface {
	GenerateToken(length int) (string, error)       // Generate URL-safe base64 token
	GenerateRandomBytes(length int) ([]byte, error) // Generate cryptographically secure random bytes
}

// CryptoService combines all cryptographic operations in a single interface
// Provides a unified API for all crypto-related functionality
type CryptoService interface {
	Hasher
	TokenGenerator
}

// cryptoService implements CryptoService with configurable bcrypt cost
type cryptoService struct {
	bcryptCost int // Cost parameter for bcrypt (higher = more secure but slower)
}

// NewCryptoService creates a new CryptoService with specified bcrypt cost
// Validates cost parameter against bcrypt min/max bounds for safety
func NewCryptoService(bcryptCost int) CryptoService {
	if bcryptCost < bcrypt.MinCost || bcryptCost > bcrypt.MaxCost {
		bcryptCost = DefaultBcryptCost
	}
	return &cryptoService{
		bcryptCost: bcryptCost,
	}
}

// NewDefaultCryptoService creates a CryptoService with recommended default settings
// Suitable for most production applications
func NewDefaultCryptoService() CryptoService {
	return NewCryptoService(DefaultBcryptCost)
}

// HashPassword generates a bcrypt hash from a plaintext password
// Returns error if password is empty or bcrypt fails
func (c *cryptoService) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), c.bcryptCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedBytes), nil
}

// CheckPasswordHash verifies a plaintext password against a bcrypt hash
// Returns false for empty inputs to prevent timing attacks
func (c *cryptoService) CheckPasswordHash(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}
	
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomBytes produces cryptographically secure random bytes
// Uses crypto/rand for CSPRNG quality randomness
func (c *cryptoService) GenerateRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be positive")
	}
	
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	
	return bytes, nil
}

// GenerateToken creates a URL-safe base64 encoded token
// Enforces minimum length for security and uses cryptographically secure randomness
func (c *cryptoService) GenerateToken(length int) (string, error) {
	if length < MinTokenBytes {
		return "", fmt.Errorf("token length must be at least %d bytes", MinTokenBytes)
	}
	
	bytes, err := c.GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SecureCompare performs constant-time string comparison to prevent timing attacks
// Essential for comparing tokens, hashes, and other security-sensitive values
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// SecureCompareBytes performs constant-time byte slice comparison
// Lower-level version of SecureCompare for byte operations
func SecureCompareBytes(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

