package repositories

import (
	"context"
	stdErrors "errors"
	"fmt"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id string, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	results := []models.User{}

	if err := u.db.WithContext(ctx).
		Select("id, name, username, email, password, role, created_at, updated_at").
		Where("deleted_at is null").Find(&results).Error; err != nil {

		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, errors.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return results, nil
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).
		Select("id, name, username, email, password, role, created_at, updated_at").
		Where("username = ?", username).First(&user).Error; err != nil {

		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return &models.User{}, errors.ErrUserNotFound
		}

		return &models.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (u *userRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := u.db.WithContext(ctx).
		Select("id, name, username, email, password, role, created_at, updated_at").
		Where("id = ?", id).First(&user).Error; err != nil {

		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return &models.User{}, errors.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.ErrInvalidQuery
	}

	res := u.db.WithContext(ctx).Create(&user)
	if res.Error != nil {
		return fmt.Errorf("failed to create user: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id string, user *models.User) error {
	result := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	result := u.db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
