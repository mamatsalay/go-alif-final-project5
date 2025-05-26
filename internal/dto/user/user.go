package user

import "workout-tracker/internal/model/user"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Username string    `json:"username"`
	Role     user.Role `json:"role"`
	ID       int       `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Username     string    `json:"username"`
	Role         user.Role `json:"role"`
	UserID       int       `json:"user_id"`
}
