package token

import (
	"context"
	"time"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models" // Jika model User digunakan di sini
)

// TokenService defines the interface for token management service (JWT).
type TokenService interface {
	// generates a JWT for the user.
	GenerateToken(ctx context.Context, user *models.User) (string, error)
	//  validates a JWT and returns user details if valid.
	ValidateToken(ctx context.Context, tokenString string) (isValid bool, userID string, username string, role string, errorMessage string, err error)
	// adds a JWT ID (JTI) to the blacklist.
	BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error
}
