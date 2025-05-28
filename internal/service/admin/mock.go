package admin

import (
	"context"
	"fmt"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/model/exercise"

	"github.com/stretchr/testify/mock"
)

type MockExerciseRepo struct {
	mock.Mock
}

func (m *MockExerciseRepo) CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), fmt.Errorf("error creating exercise: %w", args.Error(1))
}

func (m *MockExerciseRepo) UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error {
	args := m.Called(ctx, id, input)
	return fmt.Errorf("error updating exercise: %w", args.Error(0))
}

func (m *MockExerciseRepo) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	args := m.Called(ctx)

	exList, ok := args.Get(0).([]exercise.Exercise)
	if !ok {
		return nil, fmt.Errorf("invalid type for []exercise.Exercise: %w", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error in GetAllExercises: %w", err)
	}

	return exList, nil
}

func (m *MockExerciseRepo) DeleteExercise(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return fmt.Errorf("error deliting exercise: %w", args.Error(0))
}

func (m *MockExerciseRepo) GetExerciseByID(ctx context.Context, id int) (*exercise.Exercise, error) {
	args := m.Called(ctx, id)

	ex, ok := args.Get(0).(*exercise.Exercise)
	if !ok {
		return nil, fmt.Errorf("invalid type for *exercise.Exercise: %w", args.Error(1))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error in GetExerciseByID: %w", err)
	}

	return ex, nil
}
