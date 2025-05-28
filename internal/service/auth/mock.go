package auth

import (
	"context"
	"fmt"
	"time"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/user/jwt"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(ctx context.Context, u user.User) (int, error) {
	args := m.Called(ctx, u)
	err := args.Error(1)
	if err != nil {
		return 0, fmt.Errorf("error creating exercise: %w", err)
	}
	return args.Int(0), nil
}

func (m *mockUserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)

	userVal, ok := args.Get(0).(*user.User)
	if !ok {
		return nil, fmt.Errorf("invalid type for *user.User: %w", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error get user by Username: %w", err)
	}

	return userVal, nil
}

func (m *mockUserRepository) GetUserByUserID(ctx context.Context, id int) (*user.User, error) {
	args := m.Called(ctx, id)

	userVal, ok := args.Get(0).(*user.User)
	if !ok {
		return nil, fmt.Errorf("invalid type for *user.User: %w", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error get user by id: %w", err)
	}

	return userVal, nil
}

func (m *mockUserRepository) StoreRefreshToken(ctx context.Context, token string, userID int, expires time.Time) (uuid.UUID, error) {
	args := m.Called(ctx, token, userID, expires)

	tokenVal, ok := args.Get(0).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid type for uuid.:%w ", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return uuid.Nil, fmt.Errorf("error StoreRefreshToken: %w", err)
	}

	return tokenVal, nil
}

func (m *mockUserRepository) GetRefreshToken(ctx context.Context, token string) (*jwt.RefreshToken, error) {
	args := m.Called(ctx, token)

	refreshToken, ok := args.Get(0).(*jwt.RefreshToken)
	if !ok {
		return nil, fmt.Errorf("invalid type for *jwt.RefreshToken: %w", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error GetRefreshToken: %w", err)
	}

	return refreshToken, nil
}

func (m *mockUserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)

	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error DeleteRefreshToken: %w", err)
	}

	return nil
}

func (m *mockUserRepository) IncrementTokenVersion(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)

	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error IncrementTokenVersion: %w", err)
	}
	return nil
}
