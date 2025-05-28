package user

import (
	"errors"
	"testing"
	"time"
	"workout-tracker/internal/erorrs"
	"workout-tracker/internal/model/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type MockCommandTag struct{}

func (m MockCommandTag) String() string      { return "INSERT 0 1" }
func (m MockCommandTag) RowsAffected() int64 { return 1 }

func TestCreateUser_Success(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	mockRow := new(MockRow)
	log := zaptest.NewLogger(t).Sugar()

	userInput := user.User{
		Username: "newuser",
		Password: "password123",
	}
	expectedID := 42

	// Provide expected behavior
	mockPool.On("QueryRow", ctx,
		"INSERT INTO users(username, password) VALUES($1, $2) RETURNING id",
		userInput.Username, userInput.Password,
	).Return(mockRow)

	// FIX: use mock.AnythingOfType("*int")
	mockRow.On("Scan", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
		ptr, ok := args.Get(0).(*int)
		if !ok {
			t.Fatalf("Scan called with non-*int argument: %T", args.Get(0))
		}
		*ptr = expectedID
	}).Return(nil)

	repo := &UserRepository{
		Pool: mockPool,
		Log:  log,
	}

	id, err := repo.CreateUser(ctx, userInput)

	assert.NoError(t, err, "unexpected error from CreateUser")
	assert.Equal(t, expectedID, id, "returned ID should match expected")

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	mockRow := new(MockRow)
	log := zaptest.NewLogger(t).Sugar()

	userInput := user.User{
		Username: "existinguser",
		Password: "password123",
	}

	pgErr := &pgconn.PgError{Code: "23505"}

	mockPool.On("QueryRow", ctx,
		"INSERT INTO users(username, password) VALUES($1, $2) RETURNING id",
		userInput.Username, userInput.Password,
	).Return(mockRow)

	mockRow.On("Scan", mock.AnythingOfType("*int")).Return(pgErr)

	repo := &UserRepository{
		Pool: mockPool,
		Log:  log,
	}

	id, err := repo.CreateUser(ctx, userInput)

	assert.ErrorIs(t, err, erorrs.ErrUsernameAlreadyExists)
	assert.Equal(t, 0, id)

	mockPool.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	mockRow := new(MockRow)
	log := zaptest.NewLogger(t).Sugar()

	username := "nouser"

	mockPool.On("QueryRow", ctx, mock.Anything, username).Return(mockRow)
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgx.ErrNoRows)

	repo := &UserRepository{Pool: mockPool, Log: log}
	u, err := repo.GetUserByUsername(ctx, username)

	assert.Nil(t, u)
	assert.ErrorIs(t, err, erorrs.ErrUserNotFound)
}

func TestGetUserByUserID_Error(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	mockRow := new(MockRow)
	log := zaptest.NewLogger(t).Sugar()

	userID := 123
	testErr := errors.New("some db error")
	mockPool.On("QueryRow", ctx, mock.Anything, userID).Return(mockRow)
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(testErr)

	repo := &UserRepository{Pool: mockPool, Log: log}
	u, err := repo.GetUserByUserID(ctx, userID)

	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting user")
}

func TestStoreRefreshToken_Success(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	log := zaptest.NewLogger(t).Sugar()

	token := "some-token"
	userID := 123
	expires := time.Now().Add(time.Hour)

	mockPool.On("Exec", ctx, mock.Anything, mock.AnythingOfType("uuid.UUID"), userID, token, expires).
		Return(pgconn.NewCommandTag(""), nil)

	repo := &UserRepository{Pool: mockPool, Log: log}
	id, err := repo.StoreRefreshToken(ctx, token, userID, expires)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestGetRefreshToken_TokenNotFound(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	mockRow := new(MockRow)
	log := zaptest.NewLogger(t).Sugar()

	token := "missing-token"
	mockPool.On("QueryRow", ctx, mock.Anything, token).Return(mockRow)
	mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgx.ErrNoRows)

	repo := &UserRepository{Pool: mockPool, Log: log}
	rt, err := repo.GetRefreshToken(ctx, token)

	assert.Nil(t, rt)
	assert.ErrorIs(t, err, erorrs.ErrTokenNotFound)
}

func TestDeleteRefreshToken_Error(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	log := zaptest.NewLogger(t).Sugar()

	token := "tok"
	testErr := errors.New("db error")

	mockPool.On("Exec", ctx, mock.Anything, token).Return(pgconn.NewCommandTag(""), testErr)

	repo := &UserRepository{Pool: mockPool, Log: log}
	err := repo.DeleteRefreshToken(ctx, token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error deleting refresh token: db error")
}

func TestIncrementTokenVersion_Error(t *testing.T) {
	ctx := t.Context()
	mockPool := new(MockPool)
	log := zaptest.NewLogger(t).Sugar()

	userID := 321
	testErr := errors.New("db error")

	mockPool.On("Exec", ctx, mock.Anything, userID).Return(pgconn.NewCommandTag(""), testErr)

	repo := &UserRepository{Pool: mockPool, Log: log}
	err := repo.IncrementTokenVersion(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "increment token version")
}
