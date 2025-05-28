package user

import (
	"context"
	"errors"
	"fmt"
	"time"
	"workout-tracker/internal/erorrs"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/user/jwt"
	"workout-tracker/pkg/db"
	"workout-tracker/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/dig"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type UserRepositoryParams struct {
	dig.In

	DB  *db.DB
	Log logger.SugaredLoggerInterface
}

type UserRepository struct {
	Pool DBPool
	Log  logger.SugaredLoggerInterface
}

func NewRepository(params UserRepositoryParams) *UserRepository {
	return &UserRepository{
		Pool: params.DB.Pool,
		Log:  params.Log,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user user.User) (int, error) {
	var id int

	err := r.Pool.QueryRow(ctx, "INSERT INTO users(username, password) VALUES($1, $2) RETURNING id",
		user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, erorrs.ErrUsernameAlreadyExists
		}
		r.Log.Errorw("error inserting user", erorrs.ErrorKey, err, "username", user.Username)
		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	r.Log.Info("user created", "id", id)
	return id, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `SELECT id, username, password, role, createdat, updatedat, token_version FROM users WHERE username = $1`
	return r.getUser(ctx, query, username)
}

func (r *UserRepository) GetUserByUserID(ctx context.Context, id int) (*user.User, error) {
	query := `SELECT id, username, password, role, createdat, updatedat, token_version FROM users WHERE id = $1`
	return r.getUser(ctx, query, id)
}

func (r *UserRepository) getUser(ctx context.Context, query string, arg any) (*user.User, error) {
	var u user.User
	err := r.Pool.QueryRow(ctx, query, arg).
		Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.TokenVersion)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.Log.Info("user not found", "arg", arg)
			return nil, erorrs.ErrUserNotFound
		}
		r.Log.Errorw("error getting user", erorrs.ErrorKey, err, "arg", arg)
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, token string, userID int, expires time.Time) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.Pool.Exec(ctx, "INSERT INTO refresh_tokens(id, user_id, token, expires_at) VALUES($1, $2, $3, $4)",
		id, userID, token, expires)
	if err != nil {
		r.Log.Errorw("error inserting refresh token", erorrs.ErrorKey, err)
		return uuid.Nil, fmt.Errorf("error inserting refresh token: %w", err)
	}

	return id, nil
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, token string) (*jwt.RefreshToken, error) {
	var rt jwt.RefreshToken

	err := r.Pool.QueryRow(ctx, "SELECT id, user_id, token, expires_at, created_at from refresh_tokens WHERE token = $1",
		token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.Log.Info("refresh token not found", "token", token)
			return nil, erorrs.ErrTokenNotFound
		}

		r.Log.Errorw("error getting refresh token", erorrs.ErrorKey, err, "token", token)
		return nil, fmt.Errorf("error getting refresh token: %w", err)
	}

	return &rt, nil
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		r.Log.Errorw("error deleting refresh token", erorrs.ErrorKey, err)
		return fmt.Errorf("error deleting refresh token: %w", err)
	}
	return nil
}

func (r *UserRepository) IncrementTokenVersion(ctx context.Context, userID int) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE users SET token_version = token_version + 1 WHERE id = $1
	`, userID)

	if err != nil {
		r.Log.Errorw("failed to increment token version", "error", err)
		return fmt.Errorf("increment token version: %w", err)
	}
	return nil
}
