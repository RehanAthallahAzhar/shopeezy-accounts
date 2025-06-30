package models

import "time"

type Session struct {
	Token     string    `json:"token"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}
