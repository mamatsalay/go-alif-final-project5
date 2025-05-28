package user

import (
	"context"
	"time"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/user/jwt"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, u user.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByUserID(ctx context.Context, id int) (*user.User, error)
	StoreRefreshToken(ctx context.Context, token string, userID int, expires time.Time) (uuid.UUID, error)
	GetRefreshToken(ctx context.Context, token string) (*jwt.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	IncrementTokenVersion(ctx context.Context, userID int) error
}

var _ UserRepositoryInterface = (*UserRepository)(nil)
