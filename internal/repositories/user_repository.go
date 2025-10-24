package repositories

import (
	"context"
	"database/sql"
	stdErrors "errors"
	"fmt"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/db"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/entities"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entities.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, id uuid.UUID, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(sqlcQueries *db.Queries) UserRepository {
	return &userRepository{db: sqlcQueries}
}

func (u *userRepository) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	rows, err := u.db.GetAllUsers(ctx)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	users := make([]entities.User, 0, len(rows))
	for _, r := range rows {
		users = append(users, entities.User{
			ID:        r.ID,
			Name:      r.Name,
			Username:  r.Username,
			Email:     r.Email,
			Password:  r.Password,
			Role:      r.Role,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return users, nil
}

func (u *userRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	row, err := u.db.GetUserByUsername(ctx, username)
	if err != nil {
		if stdErrors.Is(err, sql.ErrNoRows) {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	user := &entities.User{
		ID:        row.ID,
		Name:      row.Name,
		Username:  row.Username,
		Email:     row.Email,
		Password:  row.Password,
		Role:      row.Role,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return user, nil
}

func (u *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	row, err := u.db.GetUserById(ctx, id)
	if err != nil {
		if stdErrors.Is(err, gorm.ErrRecordNotFound) {
			return &entities.User{}, errors.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	user := &entities.User{
		ID:        row.ID,
		Name:      row.Name,
		Username:  row.Username,
		Email:     row.Email,
		Password:  row.Password,
		Role:      row.Role,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}

	return user, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *entities.User) error {
	if user == nil {
		return errors.ErrInvalidQuery
	}

	res, err := u.db.CreateUser(ctx, db.CreateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (u *userRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *entities.User) error {

	res, err := u.db.UpdateUser(ctx, db.UpdateUserParams{
		ID:       id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {

	res, err := u.db.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}
