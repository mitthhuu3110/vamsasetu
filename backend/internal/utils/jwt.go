package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the JWT claims structure
type TokenClaims struct {
	UserID uint   `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Type   string `json:"type,omitempty"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a JWT access token with 15 minute expiry
// Parameters:
//   - userID: The unique identifier of the user
//   - email: The user's email address
//   - role: The user's role (owner, viewer, admin)
//
// Returns the signed JWT token string or an error
func GenerateAccessToken(userID uint, email, role string) (string, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Create claims with 15 minute expiry
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a JWT refresh token with 7 day expiry
// Parameters:
//   - userID: The unique identifier of the user
//
// Returns the signed JWT token string or an error
func GenerateRefreshToken(userID uint) (string, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Create claims with 7 day expiry
	now := time.Now()
	claims := TokenClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
// Parameters:
//   - tokenString: The JWT token string to validate
//
// Returns the token claims or an error if validation fails
func ValidateToken(tokenString string) (*TokenClaims, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
