package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/rehanazhar/shopeezy-account/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) ReadAllProducts(ctx context.Context) ([]models.User, error) {
	results := []models.User{}
	err := u.db.WithContext(ctx).Table("users").Select("*").Where("deleted_at is null").Find(&results).Error
	if err != nil {
		return []models.User{}, err
	}
	return results, nil
}

func (u *UserRepository) FindUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, models.ErrProductNotFound
		}
		return models.User{}, fmt.Errorf("failed to read product by ID: %w", result.Error)
	}
	return user, nil
}

func (u *UserRepository) FindUserById(ctx context.Context, id uint) (models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, models.ErrProductNotFound
		}
		return models.User{}, fmt.Errorf("failed to read product by ID: %w", result.Error)
	}
	return user, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, id uint, user *models.User) error {
	err := u.db.WithContext(ctx).Table("users").Where("id = ?", id).Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("product not found or already deleted")
	}
	return nil
}
