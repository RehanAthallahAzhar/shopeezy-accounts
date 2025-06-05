package handlers

import (
	repositories "github.com/rehanazhar/shopeezy-account/repositories"
)

// API struct yang akan memegang dependensi repository
type API struct {
	// UserRepo    repositories.UserRepository
	// SessionRepo repositories.SessionsRepository
	UserRepo repositories.UserRepository
}

// NewHandler adalah konstruktor untuk API
// Perbarui agar menerima semua repository
func NewHandler(
	// userRepo repositories.UserRepository,
	// sessionRepo repositories.SessionsRepository,
	userRepo repositories.UserRepository,
) *API {
	return &API{
		UserRepo: userRepo,
	}
}
