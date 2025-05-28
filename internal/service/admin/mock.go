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
	err := args.Error(1)
	if err != nil {
		return 0, fmt.Errorf("error creating exercise: %w", err)
	}
	return args.Int(0), nil
}

func (m *MockExerciseRepo) UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error {
	args := m.Called(ctx, id, input)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error updating exercise: %w", err)
	}
	return nil
}

func (m *MockExerciseRepo) DeleteExercise(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("error deleting exercise: %w", err)
	}
	return nil
}

func (m *MockExerciseRepo) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	args := m.Called(ctx)

	exList, ok := args.Get(0).([]exercise.Exercise)
	if !ok {
		return nil, fmt.Errorf("invalid type for []exercise.Exercise")
	}

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("error in GetAllExercises: %w", err)
	}

	return exList, nil
}

func (m *MockExerciseRepo) GetExerciseByID(ctx context.Context, id int) (*exercise.Exercise, error) {
	args := m.Called(ctx, id)

	ex, ok := args.Get(0).(*exercise.Exercise)
	if !ok {
		return nil, fmt.Errorf("invalid type for *exercise.Exercise")
	}

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("error in GetExerciseByID: %w", err)
	}

	return ex, nil
}
