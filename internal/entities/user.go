package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string         `json:"name"`
	Username    string         `json:"username" gorm:"unique;not null"`
	Email       string         `json:"email" gorm:"unique;not null"`
	Role        string         `gorm:"type:varchar(50);default:'user'"`
	Address     string         `json:"address"`
	PhoneNumber string         `json:"phone_number"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
