package admin

import (
	"context"
	"fmt"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/model/exercise"
	repo "workout-tracker/internal/repository/exercise"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AdminServiceParams struct {
	dig.In
	Log          *zap.SugaredLogger
	ExerciseRepo repo.ExerciseRepositoryInterface
}

type AdminService struct {
	Log          *zap.SugaredLogger
	ExerciseRepo repo.ExerciseRepositoryInterface
}

func NewAdminService(params AdminServiceParams) *AdminService {
	return &AdminService{
		Log:          params.Log,
		ExerciseRepo: params.ExerciseRepo,
	}
}

func (s *AdminService) CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error) {
	created, err := s.ExerciseRepo.CreateExercise(ctx, input)
	if err != nil {
		s.Log.Errorw("Service failed to create exercise", "error", err)
		return 0, fmt.Errorf("failed to create exercise in service: %w", err)
	}
	return created, nil
}

func (s *AdminService) UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error {
	err := s.ExerciseRepo.UpdateExercise(ctx, id, input)
	if err != nil {
		s.Log.Errorw("Service failed to update exercise", "error", err)
		return fmt.Errorf("failed to update exercise in service: %w", err)
	}
	return nil
}

func (s *AdminService) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	exercises, err := s.ExerciseRepo.GetAllExercises(ctx)
	if err != nil {
		s.Log.Errorw("Service failed to get exercises", "error", err)
		return nil, fmt.Errorf("failed to get exercises: %w", err)
	}
	return exercises, nil
}

func (s *AdminService) DeleteExercise(ctx context.Context, id int) error {
	err := s.ExerciseRepo.DeleteExercise(ctx, id)
	if err != nil {
		s.Log.Errorw("Service failed to delete exercise", "id", id, "error", err)
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}

func (s *AdminService) GetExerciseByID(ctx context.Context, id int) (*exercise.Exercise, error) {
	res, err := s.ExerciseRepo.GetExerciseByID(ctx, id)
	if err != nil {
		s.Log.Errorw("Service failed to get exercise", "id", id, "error", err)
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	return res, nil
}
