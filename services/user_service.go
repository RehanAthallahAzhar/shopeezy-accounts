package services

import (
	"context"
	stdErrors "errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	pkgErrors "github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services/token"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Register(ctx context.Context, userData *models.UserRegisterRequest) error
	Login(ctx context.Context, userData *models.UserLoginRequest) (*models.LoginResponse, error)
	Logout(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]models.UserResponse, error)
	GetUserById(ctx context.Context, id string) (*models.UserResponse, error)
	UpdateUser(ctx context.Context, id string, userReq *models.UserUpdateRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	userRepo         repositories.UserRepository
	validator        *validator.Validate
	tokenService     token.TokenService
	JWTBlacklistRepo repositories.JWTBlacklistRepository
}

func NewUserService(
	userRepo repositories.UserRepository,
	validator *validator.Validate,
	tokenService token.TokenService,
	JWTBlacklistRepo repositories.JWTBlacklistRepository,
) UserService {
	return &UserServiceImpl{
		userRepo:         userRepo,
		validator:        validator,
		tokenService:     tokenService,
		JWTBlacklistRepo: JWTBlacklistRepo,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, userData *models.UserRegisterRequest) error {
	if userData.Name == "" || userData.Email == "" || userData.Password == "" {
		return pkgErrors.ErrInvalidRequestPayload
	}

	if userData.Token != "secret token from HRD or Manager" {
		return pkgErrors.ErrInvalidTokenRole
	}

	if userData.Role == "" {
		userData.Role = "user"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("service: Failed to register user: %w", err)
	}

	userData.Password = string(hashedPassword)

	user := &models.User{
		Name:     userData.Name,
		Username: userData.Username,
		Email:    userData.Email,
		Password: userData.Password,
		Role:     userData.Role,
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return fmt.Errorf("service: failed to register user: %w", err)
	}

	return nil
}

func (s *UserServiceImpl) Login(ctx context.Context, userData *models.UserLoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, userData.Username)
	if err != nil {
		if stdErrors.Is(err, pkgErrors.ErrUserNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("service: failed to login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, pkgErrors.ErrInvalidCredentials
	}

	signedToken, err := s.tokenService.GenerateToken(ctx, user)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return nil, fmt.Errorf("service: Failed to generate JWT: %w", err)
	}

	res := &models.LoginResponse{
		Id:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Token:     signedToken,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	return res, nil
}

func (s *UserServiceImpl) Logout(ctx context.Context, authHeader string) error {

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return pkgErrors.ErrInvalidTokenFormat
	}

	/*
		Parse JWT to get claims, specifically JTI and expiration time
		We use ParseUnverified because the signature will be verified by the gRPC AuthServer
		For logout, we only need to read the claims for JTI and expiration time.
	*/

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &models.JWTClaims{})
	if err != nil {
		log.Printf("Error parsing unverified JWT for logout: %v", err)
		return pkgErrors.ErrInvalidToken
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok {
		return pkgErrors.ErrInvalidToken
	}

	jti := claims.ID
	if jti == "" {
		return pkgErrors.ErrMissingJTI
	}

	// Calculate remaining time the token is valid
	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime < 0 {
		return pkgErrors.ErrExpiredToken // Token is already expired, no need to blacklist
	}

	// Add JTI to Redis blacklist using JWTBlacklistRepo
	err = s.JWTBlacklistRepo.AddToBlacklist(ctx, jti, remainingTime)
	if err != nil {
		log.Printf("Error adding JTI %s to blacklist: %v", jti, err)
		return pkgErrors.ErrFailedToRevokeToken
	}

	log.Printf("Token with JTI %s successfully revoked.", jti)
	return nil
}

func (s *UserServiceImpl) GetAllUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		if stdErrors.Is(err, pkgErrors.ErrUserNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("service: failed to get all users: %w", err)
	}

	var res []models.UserResponse
	for _, user := range users {
		res = append(res, *mapToUserResponse(&user))
	}

	return res, nil
}

func (s *UserServiceImpl) GetUserById(ctx context.Context, id string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		if stdErrors.Is(err, pkgErrors.ErrUserNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("service: failed to get user by id: %w", err)
	}

	return mapToUserResponse(user), nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, userReq *models.UserUpdateRequest) (*models.UserResponse, error) {
	if err := s.validator.Struct(userReq); err != nil {
		return nil, fmt.Errorf("%w: %s", pkgErrors.ErrInvalidRequestPayload, err)
	}

	user := &models.User{
		Name:     userReq.Name,
		Username: userReq.Username,
		Email:    userReq.Email,
	}

	if userReq.Password != "" {
		user.Password = userReq.Password
	}

	err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		if stdErrors.Is(err, pkgErrors.ErrUserNotFound) {
			return nil, pkgErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("UpdateUser service error: %w", err)
	}

	updatedUser, err := s.userRepo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapToUserResponse(updatedUser), nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		if stdErrors.Is(err, pkgErrors.ErrUserNotFound) {
			return pkgErrors.ErrUserNotFound
		}
		return fmt.Errorf("UpdateUser service error: %w", err)
	}

	return nil
}

func mapToUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
