package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `json:"name"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"password"`
	Role      string         `gorm:"type:varchar(50);default:'user'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserAuthRequest struct {
	Name      string `json:"name" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role" validate:"required"`
	TokenRole string `json:"token"`
}

type UserResponse struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Token     *string `json:"token,omitempty"` // omitempty = toeken akan ditampilkan jika bernilai tidak nil
}

// Struct ini digunakan untuk menerima data saat proses create user
type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required"`
	Role     string `gorm:"type:varchar(50);default:'user'"`
}

// Struct ini digunakan untuk menerima data saat proses update user
type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Email    string `json:"email" binding:"required" gorm:"unique;not null"`
	Password string `json:"password,omitempty"`
}

// Struct ini digunakan saat user melakukan proses login
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
