package auth

import (
	"context"
	"fmt"
	"os"
	"time"
	"workout-tracker/internal/erorrs"
	model "workout-tracker/internal/model/user"
	"workout-tracker/internal/repository/user"

	"golang.org/x/crypto/bcrypt"

	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const oneMonth = time.Hour * 24 * 30
const halfAnHour = time.Minute * 30

type AuthServiceParams struct {
	dig.In

	Repo *user.UserRepository
	Log  *zap.SugaredLogger
}

type AuthService struct {
	Repo   *user.UserRepository
	Log    *zap.SugaredLogger
	Secret string
}

func NewAuthService(params AuthServiceParams) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	return &AuthService{
		Repo:   params.Repo,
		Log:    params.Log,
		Secret: secret,
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Errorw("failed to hash password", erorrs.ErrorKey, err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), err
}

func (s *AuthService) CheckPassword(hashed, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		s.Log.Errorw("failed to compare password", erorrs.ErrorKey, err)
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int, error) {
	res, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		s.Log.Errorw("failed to create user", erorrs.ErrorKey, err)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return res, nil
}

func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.Repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.Log.Errorw("failed to get user by username", "username", username, erorrs.ErrorKey, err)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (s *AuthService) GenerateAccessToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(halfAnHour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	res, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		s.Log.Errorw("error generating access token", erorrs.ErrorKey, err)
		return "", fmt.Errorf("error generating access token: %w", err)
	}

	return res, nil
}

func (s *AuthService) GenerateAndStoreRefreshToken(ctx context.Context, userID int) (string, error) {
	token := uuid.NewString()
	expires := time.Now().Add(oneMonth)

	_, err := s.Repo.StoreRefreshToken(ctx, token, userID, expires)
	if err != nil {
		s.Log.Errorw("failed to store refresh token", erorrs.ErrorKey, err)
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}
