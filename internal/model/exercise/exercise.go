package exercise

import "time"

type Exercise struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"update_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Name        string     `json:"name"`
	ID          int        `json:"id"`
	Description string     `json:"description"`
}
