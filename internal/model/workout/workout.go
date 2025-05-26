package workout

import (
	"time"
)

type Workout struct {
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	Category  string     `json:"category"`
	Title     string     `json:"title"`
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
}
