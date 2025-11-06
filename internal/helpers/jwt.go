package helpers

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTHelper is a struct that holds the secret key for generating tokens.
type JWTHelper struct {
	secretKey []byte
}

// NewJWTHelper creates a new instance of JWTHelper.
// The secret key is passed in here during initialization.
func NewJWTHelper(secret string) *JWTHelper {
	return &JWTHelper{
		secretKey: []byte(secret),
	}
}

// GenerateToken creates a new JWT for a given username.
func (h *JWTHelper) GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key stored in the struct
	signedToken, err := token.SignedString(h.secretKey)
	if err != nil {
		return "", err // It's better to return the error
	}

	return signedToken, nil
}

// You can also add your token validation function here as another method.
// func (h *JWTHelper) ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) { ... }
