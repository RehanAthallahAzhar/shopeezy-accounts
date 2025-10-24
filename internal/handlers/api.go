package handlers

import (
	"log"

	stdErr "errors"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services/token"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type API struct {
	UserRepo         repositories.UserRepository
	UserService      services.UserService
	TokenService     token.TokenService
	JWTBlacklistRepo repositories.JWTBlacklistRepository
}

func NewHandler(
	userRepo repositories.UserRepository,
	userService services.UserService,
	tokenService token.TokenService,
	jwtBlacklistRepo repositories.JWTBlacklistRepository) *API {
	return &API{
		UserRepo:         userRepo,
		UserService:      userService,
		TokenService:     tokenService,
		JWTBlacklistRepo: jwtBlacklistRepo,
	}
}

func extractUserID(c echo.Context) (uuid.UUID, error) {
	val := c.Get("userID")

	if val == nil {
		log.Println("DEBUG: c.Get('userID') hasilnya nil.")
		return uuid.Nil, stdErr.New("invalid user session: userID is nil in context") // Gunakan error yang lebih jelas
	}

	if id, ok := val.(uuid.UUID); ok {
		return id, nil
	}

	log.Printf("DEBUG: Gagal assertion ke uuid.UUID. Tipe data aktual adalah %T", val)
	return uuid.Nil, stdErr.New("invalid user session: userID in context is not of type uuid.UUID")
}
