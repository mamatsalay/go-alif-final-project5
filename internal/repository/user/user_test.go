package user

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	errs "workout-tracker/internal/erorrs"
	"workout-tracker/internal/model/user"
	"workout-tracker/pkg/db"
)

var userRepo *UserRepository

func TestMain(m *testing.M) {
	t := &testing.T{}
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")

	pooler, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	resource, err := pooler.Run("postgres", "13-alpine", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		resource.GetPort("5432/tcp"), os.Getenv("DB_NAME"),
	)

	pooler.MaxWait = 30 * time.Second
	t.Setenv("DB_URL", dsn)
	err = pooler.Retry(func() error {
		database, err := db.New(zap.NewNop().Sugar())
		if err != nil {
			return err
		}
		return database.Pool.Ping(context.Background())
	})
	if err != nil {
		panic(err)
	}

	database, err := db.New(zap.NewNop().Sugar())
	if err != nil {
		panic(err)
	}

	exec := func(q string) {
		if _, err := database.Pool.Exec(context.Background(), q); err != nil {
			panic(err)
		}
	}
	exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'user',
		createdat TIMESTAMP DEFAULT NOW(),
		updatedat TIMESTAMP DEFAULT NOW(),
		token_version INT DEFAULT 0
	);`)
	exec(`CREATE TABLE refresh_tokens (
		id UUID PRIMARY KEY,
		user_id INT REFERENCES users(id),
		token TEXT UNIQUE NOT NULL,
		expires_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT NOW()
	);`)

	userRepo = NewRepository(UserRepositoryParams{
		DB:  database,
		Log: zap.NewNop().Sugar(),
	})

	code := m.Run()
	database.Pool.Close()
	_ = pooler.Purge(resource)
	os.Exit(code)
}

func TestCreateAndGetUser(t *testing.T) {
	ctx := context.Background()
	u := user.User{Username: "alice", Password: "pass"}

	id, err := userRepo.CreateUser(ctx, u)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := userRepo.GetUserByUsername(ctx, u.Username)
	require.NoError(t, err)
	assert.Equal(t, id, got.ID)
	assert.Equal(t, u.Username, got.Username)

	byID, err := userRepo.GetUserByUserID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, u.Username, byID.Username)
}

func TestCreateUser_Duplicate(t *testing.T) {
	ctx := context.Background()
	u := user.User{Username: "bob", Password: "pw"}
	_, err := userRepo.CreateUser(ctx, u)
	require.NoError(t, err)

	_, err = userRepo.CreateUser(ctx, u)
	assert.ErrorIs(t, err, errs.ErrUsernameAlreadyExists)
}

func TestRefreshTokenCRUD(t *testing.T) {
	ctx := context.Background()
	u := user.User{Username: "charlie", Password: "pw2"}
	uid, err := userRepo.CreateUser(ctx, u)
	require.NoError(t, err)

	expires := time.Now().Add(time.Hour)
	tok := "token123"
	rtID, err := userRepo.StoreRefreshToken(ctx, tok, uid, expires)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, rtID)

	rt, err := userRepo.GetRefreshToken(ctx, tok)
	require.NoError(t, err)
	assert.Equal(t, tok, rt.Token)
	assert.Equal(t, uid, rt.UserID)

	err = userRepo.DeleteRefreshToken(ctx, tok)
	require.NoError(t, err)

	_, err = userRepo.GetRefreshToken(ctx, tok)
	assert.ErrorIs(t, err, errs.ErrTokenNotFound)
}

func TestIncrementTokenVersion(t *testing.T) {
	ctx := context.Background()
	u := user.User{Username: "dave", Password: "pw3"}
	uid, err := userRepo.CreateUser(ctx, u)
	require.NoError(t, err)

	orig, err := userRepo.GetUserByUserID(ctx, uid)
	require.NoError(t, err)
	assert.Equal(t, 0, orig.TokenVersion)

	err = userRepo.IncrementTokenVersion(ctx, uid)
	require.NoError(t, err)

	upd, err := userRepo.GetUserByUserID(ctx, uid)
	require.NoError(t, err)
	assert.Equal(t, orig.TokenVersion+1, upd.TokenVersion)
}
