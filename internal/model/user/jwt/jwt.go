package jwt

import (
	"time"
)

type RefreshToken struct {
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token"`
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
}
