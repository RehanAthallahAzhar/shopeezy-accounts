package handlers

import (
	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services"
)

type API struct {
	UserRepo repositories.UserRepository
	TokenSvc services.TokenService
}

func NewHandler(userRepo repositories.UserRepository, tokenSvc services.TokenService) *API {
	return &API{
		UserRepo: userRepo,
		TokenSvc: tokenSvc,
	}
}

// LoginUser menangani permintaan login pengguna
