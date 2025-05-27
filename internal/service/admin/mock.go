package admin

import (
	"context"
	"github.com/stretchr/testify/mock"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/model/exercise"
)

type MockExerciseRepo struct {
	mock.Mock
}

func (m *MockExerciseRepo) CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error) {
	args := m.Called(ctx, input)
	return args.Int(0), args.Error(1)
}

func (m *MockExerciseRepo) UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error {
	args := m.Called(ctx, id, input)
	return args.Error(0)
}

func (m *MockExerciseRepo) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	args := m.Called(ctx)
	return args.Get(0).([]exercise.Exercise), args.Error(1)
}

func (m *MockExerciseRepo) DeleteExercise(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockExerciseRepo) GetExerciseByID(ctx context.Context, id int) (*exercise.Exercise, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*exercise.Exercise), args.Error(1)
}
