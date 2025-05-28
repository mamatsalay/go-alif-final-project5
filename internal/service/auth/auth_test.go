package auth

import (
	"errors"
	"os"
	"testing"
	"time"
	"workout-tracker/internal/erorrs"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/user/jwt"

	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestEnvironment() {
	err := os.Setenv("JWT_SECRET", "testsecret")
	if err != nil {
		panic(err)
	}
}

func getTestLogger() *zap.SugaredLogger {
	logger := zap.NewNop()
	return logger.Sugar()
}

func TestAuthService_CreateUser(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{
		Repo: repo,
		Log:  getTestLogger(),
	})

	ctx := t.Context()
	userInput := user.User{
		Username: "testuser",
		Password: "hashed",
		Role:     user.UserRole,
		ID:       1,
	}

	repo.On("CreateUser", ctx, userInput).Return(1, nil)

	id, err := service.CreateUser(ctx, userInput)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	repo.AssertExpectations(t)
}

func TestAuthService_CreateUser_Error(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{
		Repo: repo,
		Log:  getTestLogger(),
	})

	ctx := t.Context()
	userInput := user.User{
		Username: "testuser",
		Password: "hashed",
		Role:     user.UserRole,
		ID:       1,
	}

	expectedError := erorrs.ErrUsernameAlreadyExists
	repo.On("CreateUser", ctx, userInput).Return(0, expectedError)

	id, err := service.CreateUser(ctx, userInput)
	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "failed to create user")
	repo.AssertExpectations(t)
}

func TestAuthService_GenerateAccessToken(t *testing.T) {
	setupTestEnvironment()

	service := NewAuthService(AuthServiceParams{Repo: nil, Log: getTestLogger()})

	usr := &user.User{
		ID:           123,
		Role:         user.UserRole,
		TokenVersion: 1,
	}

	token, err := service.GenerateAccessToken(usr)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_HashAndCheckPassword(t *testing.T) {
	setupTestEnvironment()

	service := NewAuthService(AuthServiceParams{Log: getTestLogger(), Repo: nil})

	pass := "123456"
	hash, err := service.HashPassword(pass)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = service.CheckPassword(hash, pass)
	assert.NoError(t, err)

	err = service.CheckPassword(hash, "wrong")
	assert.Error(t, err)
}

func TestAuthService_GetUserByUsername(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	expectedUser := &user.User{
		ID:       1,
		Username: "testuser",
		Role:     user.UserRole,
	}

	repo.On("GetUserByUsername", ctx, "testuser").Return(expectedUser, nil)

	result, err := service.GetUserByUsername(ctx, "testuser")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	repo.AssertExpectations(t)
}

func TestAuthService_GetUserByUsername_Error(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	repo.On("GetUserByUsername", ctx, "ghost").Return(nil, erorrs.ErrUserNotFound)

	_, err := service.GetUserByUsername(ctx, "ghost")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user by username")
	repo.AssertExpectations(t)
}

func TestAuthService_GetUserByUserID(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	expectedUser := &user.User{
		ID:       123,
		Username: "testuser",
		Role:     user.UserRole,
	}

	repo.On("GetUserByUserID", ctx, 123).Return(expectedUser, nil)

	result, err := service.GetUserByUserID(ctx, 123)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	repo.AssertExpectations(t)
}

func TestAuthService_GetUserByUserID_Error(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	repo.On("GetUserByUserID", ctx, 404).Return(nil, erorrs.ErrUserNotFound)

	_, err := service.GetUserByUserID(ctx, 404)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user by user id")
	repo.AssertExpectations(t)
}

func TestAuthService_GenerateAndStoreRefreshToken(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()
	userID := 1

	expectedUUID := uuid.New()
	repo.On("StoreRefreshToken", ctx, mock.AnythingOfType("string"), userID, mock.AnythingOfType("time.Time")).Return(expectedUUID, nil)

	token, err := service.GenerateAndStoreRefreshToken(ctx, userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	usr := &user.User{
		ID:           1,
		Role:         user.UserRole,
		TokenVersion: 2,
	}

	repo.On("IncrementTokenVersion", ctx, usr.ID).Return(nil)
	repo.On("GetUserByUserID", ctx, usr.ID).Return(usr, nil)

	token, err := service.RefreshAccessToken(ctx, usr.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestAuthService_RefreshAccessToken_IncrementTokenVersionError(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: getTestLogger()})
	ctx := t.Context()

	userID := 1
	expectedError := errors.New("database error")

	repo.On("IncrementTokenVersion", ctx, userID).Return(expectedError)

	token, err := service.RefreshAccessToken(ctx, userID)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "failed to invalidate old tokens")
	repo.AssertExpectations(t)
}

func TestAuthService_UpdateRefreshToken(t *testing.T) {
	setupTestEnvironment()

	repo := new(mockUserRepository)
	service := NewAuthService(AuthServiceParams{Repo: repo, Log: nil})
	ctx := t.Context()

	refreshToken := "valid-refresh-token"
	tokenID := uuid.New()
	refreshTokenData := &jwt.RefreshToken{
		ID:        tokenID.String(),
		UserID:    1,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	user := &user.User{
		ID:           1,
		Role:         user.UserRole,
		TokenVersion: 2,
	}

	repo.On("GetRefreshToken", ctx, refreshToken).Return(refreshTokenData, nil)
	repo.On("DeleteRefreshToken", ctx, refreshToken).Return(nil)
	repo.On("IncrementTokenVersion", ctx, refreshTokenData.UserID).Return(nil)
	repo.On("GetUserByUserID", ctx, refreshTokenData.UserID).Return(user, nil)
	newTokenID := uuid.New()
	repo.On("StoreRefreshToken", ctx, mock.AnythingOfType("string"), refreshTokenData.UserID, mock.AnythingOfType("time.Time")).Return(newTokenID, nil)

	accessToken, newRefreshToken, err := service.UpdateRefreshToken(ctx, refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, newRefreshToken)
	repo.AssertExpectations(t)
}

func TestAuthService_NewAuthService_PanicOnMissingSecret(t *testing.T) {
	err := os.Unsetenv("JWT_SECRET")
	if err != nil {
		panic(err)
	}

	assert.Panics(t, func() {
		NewAuthService(AuthServiceParams{Repo: nil, Log: getTestLogger()})
	})
}
