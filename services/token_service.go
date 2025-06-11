package services

import (
	"context"
	"fmt"
	"os" // Diperlukan jika Anda menggunakan os.Getenv untuk secretKey
	"time"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/golang-jwt/jwt/v5"
)

// TokenService mendefinisikan interface untuk layanan token.
type TokenService interface {
	// KOREKSI: Tambahkan 'role' sebagai nilai kembali di Validate
	Validate(ctx context.Context, tokenString string) (isValid bool, userID string, username string, role string, errorMessage string, err error)
	// KOREKSI: Tambahkan 'role' sebagai parameter di GenerateToken
	GenerateToken(userID string, username string, role string) (string, error)
}

// NewTokenService membuat instance TokenService baru.
func NewTokenService(userRepo repositories.UserRepository) TokenService {
	return &tokenServiceImpl{
		userRepo:  userRepo,
		secretKey: os.Getenv("JWT_SECRET"), // Pastikan JWT_SECRET ada di .env Anda!
	}
}

type tokenServiceImpl struct {
	userRepo  repositories.UserRepository
	secretKey string
}

// Validate mengimplementasikan logika validasi token.
// KOREKSI: Signature fungsi diubah untuk mengembalikan 'role'
func (s *tokenServiceImpl) Validate(ctx context.Context, tokenString string) (bool, string, string, string, string, error) {
	if tokenString == "" {
		return false, "", "", "", "Token tidak boleh kosong", nil // <-- Tambahkan string kosong untuk role
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		// var validationErr *jwt.ValidationError
		// if errors.As(err, &validationErr) {
		// 	if validationErr.Errors == jwt.ValidationErrorExpired {
		// 		return false, "", "", "", "Token kedaluwarsa", jwt.ErrTokenExpired // <-- Tambahkan string kosong untuk role
		// 	}
		// }
		return false, "", "", "", "Token tidak valid", fmt.Errorf("failed to parse token: %w", err) // <-- Tambahkan string kosong untuk role
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, "", "", "", "Token tidak valid", nil // <-- Tambahkan string kosong untuk role
	}

	userIDRaw, ok := claims["user_id"]
	if !ok {
		return false, "", "", "", "Klaim 'user_id' tidak ditemukan", fmt.Errorf("claim 'user_id' not found in token") // <-- Tambahkan string kosong untuk role
	}
	userID, ok := userIDRaw.(string)
	if !ok {
		if floatID, isFloat := userIDRaw.(float64); isFloat {
			userID = fmt.Sprintf("%.0f", floatID)
		} else {
			return false, "", "", "", "Format 'user_id' tidak valid", fmt.Errorf("invalid type for claim 'user_id': %T", userIDRaw) // <-- Tambahkan string kosong untuk role
		}
	}

	usernameRaw, ok := claims["sub"]
	if !ok {
		return false, "", "", "", "Klaim 'sub' (username) tidak ditemukan", fmt.Errorf("claim 'sub' (username) not found in token") // <-- Tambahkan string kosong untuk role
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return false, "", "", "", "Format 'username' (sub) tidak valid", fmt.Errorf("invalid type for claim 'sub' (username): %T", usernameRaw) // <-- Tambahkan string kosong untuk role
	}

	// KOREKSI: Ambil klaim 'role'
	roleRaw, ok := claims["role"]
	if !ok {
		return false, "", "", "", "Klaim 'role' tidak ditemukan", fmt.Errorf("claim 'role' not found in token") // <-- Tambahkan string kosong untuk role
	}
	userRole, ok := roleRaw.(string)
	if !ok {
		return false, "", "", "", "Format 'role' tidak valid", fmt.Errorf("invalid type for claim 'role': %T", roleRaw) // <-- Tambahkan string kosong untuk role
	}

	_, err = s.userRepo.FindUserById(ctx, userID) // Sesuaikan dengan method repo Anda
	if err != nil {
		return false, "", "", "", "User tidak ditemukan atau tidak aktif", fmt.Errorf("user not found or inactive: %v", err) // <-- Tambahkan string kosong untuk role
	}

	return true, userID, username, userRole, "", nil // <-- KOREKSI: Kembalikan userRole
}

// GenerateToken membuat JWT baru untuk pengguna yang diberikan.
func (s *tokenServiceImpl) GenerateToken(userID string, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"sub":     username,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
