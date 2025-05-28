package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	exerciseDTO "workout-tracker/internal/dto/exercise"
	workoutDTO "workout-tracker/internal/dto/workout"
	"workout-tracker/internal/handler"
	"workout-tracker/internal/handler/admin"
	"workout-tracker/internal/handler/auth"
	"workout-tracker/internal/handler/workout"
	"workout-tracker/internal/model/exercise"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/model/workoutexercisejoin"
)

type mockAuthService struct{}

func (m *mockAuthService) HashPassword(password string) (string, error) {
	return "hashed", nil
}
func (m *mockAuthService) CreateUser(ctx context.Context, u user.User) (int, error) {
	return 1, nil
}
func (m *mockAuthService) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	return &user.User{Username: username, ID: 1, TokenVersion: 1, Role: user.UserRole}, nil
}
func (m *mockAuthService) CheckPassword(hashed, password string) error {
	return nil
}
func (m *mockAuthService) GenerateAccessToken(u *user.User) (string, error) {
	return "access", nil
}
func (m *mockAuthService) GenerateAndStoreRefreshToken(ctx context.Context, userID int) (string, error) {
	return "refresh", nil
}
func (m *mockAuthService) UpdateRefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "newAccess", "newRefresh", nil
}
func (m *mockAuthService) GetUserByUserID(ctx context.Context, id int) (*user.User, error) {
	return &user.User{ID: id, TokenVersion: 1, Role: user.UserRole}, nil
}

type mockAdminService struct{}

func (m *mockAdminService) CreateExercise(ctx context.Context, req exerciseDTO.CreateExerciseRequest) (int, error) {
	return 1, nil
}
func (m *mockAdminService) UpdateExercise(ctx context.Context, id int, req exerciseDTO.CreateExerciseRequest) error {
	return nil
}
func (m *mockAdminService) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	return []exercise.Exercise{}, nil
}
func (m *mockAdminService) DeleteExercise(ctx context.Context, id int) error {
	return nil
}

type mockWorkoutService struct{}

func (m *mockWorkoutService) CreateWorkout(
	ctx context.Context,
	userID int,
	name,
	title,
	category string,
	exercises []workoutexercisejoin.WorkoutExercise) error {
	return nil
}
func (m *mockWorkoutService) UpdateWorkout(
	ctx context.Context,
	userID,
	workoutID int,
	name,
	title,
	category string,
	exercises []workoutexercisejoin.WorkoutExercise) error {
	return nil
}
func (m *mockWorkoutService) DeleteWorkout(ctx context.Context, userID, workoutID int) error {
	return nil
}
func (m *mockWorkoutService) GetAllWorkoutsWithExercises(ctx context.Context, userID int) ([]workoutDTO.WorkoutWithExercises, error) {
	return []workoutDTO.WorkoutWithExercises{}, nil
}
func (m *mockWorkoutService) GetWorkoutByID(ctx context.Context, userID, workoutID int) (*workoutDTO.WorkoutWithExercises, error) {
	return &workoutDTO.WorkoutWithExercises{}, nil
}

func TestSetupRoutes(t *testing.T) {
	t.Setenv("JWT_SECRET", "testsecret")
	defer func() {
		err := os.Unsetenv("JWT_SECRET")
		if err != nil {
			t.Fatal(err)
		}
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger := zap.NewNop().Sugar()

	authHandler := auth.NewAuthHandler(auth.AuthHandlerParams{
		Service: &mockAuthService{},
		Logger:  logger,
	})
	adminHandler := admin.NewAdminHandler(admin.AdminHandlerParams{
		Service: &mockAdminService{},
		Logger:  logger,
	})
	workoutHandler := workout.NewWorkoutHandler(workout.WorkoutHandlerParams{
		Service: &mockWorkoutService{},
		Logger:  logger,
	})

	mw := handler.NewMiddleware(handler.MiddlewareParams{
		Log:     logger,
		Service: &mockAuthService{},
	})

	SetupRoutes(router, authHandler, adminHandler, workoutHandler, mw)

	req, _ := http.NewRequest(http.MethodGet, "/workouts", http.NoBody)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
		"role":    string(user.UserRole),
		"version": float64(1),
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute)).Unix(),
	})
	tokStr, _ := token.SignedString([]byte("testsecret"))
	req.Header.Set("Authorization", "Bearer "+tokStr)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusUnauthorized, resp.Code, "Expected /workouts to be accessible with valid token")
	assert.NotEqual(t, http.StatusNotFound, resp.Code, "Route /workouts should exist")
}
