package admin

import (
	"context"
	"fmt"
	"os"
	"workout-tracker/internal/model/exercise"
	repo "workout-tracker/internal/repository/exercise"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AdminServiceParams struct {
	dig.In
	Log          *zap.SugaredLogger
	ExerciseRepo *repo.ExerciseRepository
}

type AdminService struct {
	Log          *zap.SugaredLogger
	ExerciseRepo *repo.ExerciseRepository
}

func NewAdminService(params AdminServiceParams) *AdminService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	return &AdminService{
		Log:          params.Log,
		ExerciseRepo: params.ExerciseRepo,
	}
}

// CreateExercise создаёт занятие.
func (s *AdminService) CreateExercise(ctx context.Context, ex exercise.Exercise) (exercise.Exercise, error) {
	created, err := s.ExerciseRepo.CreateExercise(ctx, ex)
	if err != nil {
		s.Log.Errorw("Service failed to create exercise", "error", err)
		return exercise.Exercise{}, fmt.Errorf("failed to create exercise in service: %w", err)
	}
	return created, nil
}

// UpdateExercise обновляет занятие.
func (s *AdminService) UpdateExercise(ctx context.Context, ex exercise.Exercise) (exercise.Exercise, error) {
	updated, err := s.ExerciseRepo.UpdateExercise(ctx, ex)
	if err != nil {
		s.Log.Errorw("Service failed to update exercise", "id", ex.ID, "error", err)
		return exercise.Exercise{}, fmt.Errorf("failed to update exercise: %w", err)
	}
	return updated, nil
}

// GetAllExercises выводит занятия.
func (s *AdminService) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	exercises, err := s.ExerciseRepo.GetAllExercises(ctx)
	if err != nil {
		s.Log.Errorw("Service failed to get exercises", "error", err)
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}
	return exercises, nil
}
