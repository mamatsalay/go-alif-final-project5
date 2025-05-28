package auth

import (
	"context"
	model "workout-tracker/internal/model/user"
	"workout-tracker/internal/service/auth"
)

type AuthServiceInterface interface {
	HashPassword(password string) (string, error)
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	CheckPassword(hashed, password string) error
	GenerateAccessToken(user *model.User) (string, error)
	GenerateAndStoreRefreshToken(ctx context.Context, userID int) (string, error)
	UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

var _ AuthServiceInterface = (*auth.AuthService)(nil)
