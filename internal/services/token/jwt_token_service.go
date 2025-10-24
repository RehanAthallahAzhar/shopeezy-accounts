package token

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/entities"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/repositories"
)

// Custom JWT Claims (must be consistent across the application)
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

type jwtTokenService struct {
	jwtSecret        string
	jwtBlacklistRepo repositories.JWTBlacklistRepository
}

// NewJWTTokenService creates a new JWTTokenService instance.
func NewJWTTokenService(jwtSecret string, jwtBlacklistRepo repositories.JWTBlacklistRepository) TokenService {
	return &jwtTokenService{
		jwtSecret:        jwtSecret,
		jwtBlacklistRepo: jwtBlacklistRepo,
	}
}

func (s *jwtTokenService) GenerateToken(ctx context.Context, user *entities.User) (string, error) {
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(), // Unique JTI (JWT ID) for blacklisting
			Issuer:    "shopeezy-account-service",
			Subject:   user.Username,
			Audience:  jwt.ClaimStrings{"shopeezy-cashier-app"}, // Audience if needed
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return signedToken, nil
}

func (s *jwtTokenService) ValidateToken(ctx context.Context, tokenString string) (isValid bool, userID uuid.UUID, username string, role string, errorMessage string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		log.Printf("Failed to parse or validate JWT: %v", err)
		return false, uuid.Nil, "", "", "Token invalid or expired", err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Printf("Invalid token or claims mismatch: %v", err)
		return false, uuid.Nil, "", "", "Invalid token", nil // No Go error, just invalid token
	}

	// Check Redis Blacklist (if JTI exists)
	jti := claims.ID // Unique JWT ID
	if jti != "" {   // JTI might be empty if not set during token creation
		isBlacklisted, err := s.jwtBlacklistRepo.IsBlacklisted(ctx, jti)
		if err != nil {
			log.Printf("Error checking JWT blacklist for JTI %s: %v", jti, err)
			return false, uuid.Nil, "", "", "Internal server error during token validation", err
		}
		if isBlacklisted {
			log.Printf("Token with JTI %s is blacklisted.", jti)
			return false, uuid.Nil, "", "", "Token has been revoked", nil // No Go error, just invalid token
		}
	}

	// Token is valid and not blacklisted
	return true, claims.UserID, claims.Username, claims.Role, "", nil
}

// BlacklistToken adds a JWT ID (JTI) to the blacklist.
func (s *jwtTokenService) BlacklistToken(ctx context.Context, jti string, expiration time.Duration) error {
	return s.jwtBlacklistRepo.AddToBlacklist(ctx, jti, expiration)
}
