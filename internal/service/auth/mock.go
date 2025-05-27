package auth

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/user/jwt"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(ctx context.Context, u user.User) (int, error) {
	args := m.Called(ctx, u)
	return args.Int(0), args.Error(1)
}

func (m *mockUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByUserID(ctx context.Context, id int) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *mockUserRepository) StoreRefreshToken(ctx context.Context, token string, userID int, expires time.Time) (uuid.UUID, error) {
	args := m.Called(ctx, token, userID, expires)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockUserRepository) GetRefreshToken(ctx context.Context, token string) (*jwt.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.RefreshToken), args.Error(1)
}

func (m *mockUserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockUserRepository) IncrementTokenVersion(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
