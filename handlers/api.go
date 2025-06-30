package handlers

import (
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services/token"
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

func extractUserID(c echo.Context) (string, error) {
	if val := c.Get("userID"); val != nil {
		if id, ok := val.(string); ok {
			return id, nil
		}
	}
	return "", errors.ErrInvalidUserSession
}
