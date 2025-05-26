package user

import "time"

type Role string

const AdminRole = Role("admin")
const UserRole = Role("user")

type User struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	Role         Role      `json:"role"`
	ID           int       `json:"id"`
	TokenVersion int       `json:"token_version"`
}
