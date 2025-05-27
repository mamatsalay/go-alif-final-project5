package workout

import (
	"context"

	"workout-tracker/internal/model/workout"
	"workout-tracker/internal/model/workoutexercisejoin"

	"github.com/stretchr/testify/mock"
)

type WorkoutRepoMock struct {
	mock.Mock
}

func (m *WorkoutRepoMock) CreateWorkout(ctx context.Context, w workout.Workout) (int, error) {
	args := m.Called(ctx, w)
	return args.Int(0), args.Error(1)
}

func (m *WorkoutRepoMock) BulkInsertWorkoutExercises(ctx context.Context, ex []workoutexercisejoin.WorkoutExercise) error {
	args := m.Called(ctx, ex)
	return args.Error(0)
}

func (m *WorkoutRepoMock) UpdateWorkout(ctx context.Context, w workout.Workout) error {
	args := m.Called(ctx, w)
	return args.Error(0)
}

func (m *WorkoutRepoMock) DeleteWorkoutExercises(ctx context.Context, workoutID int) error {
	args := m.Called(ctx, workoutID)
	return args.Error(0)
}

func (m *WorkoutRepoMock) DeleteWorkout(ctx context.Context, workoutID, userID int) error {
	args := m.Called(ctx, workoutID, userID)
	return args.Error(0)
}

func (m *WorkoutRepoMock) GetAllWorkouts(ctx context.Context, userID int) ([]workout.Workout, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]workout.Workout), args.Error(1)
}

func (m *WorkoutRepoMock) GetWorkoutByID(ctx context.Context, workoutID, userID int) (*workout.Workout, error) {
	args := m.Called(ctx, workoutID, userID)
	return args.Get(0).(*workout.Workout), args.Error(1)
}

func (m *WorkoutRepoMock) GetWorkoutExercises(ctx context.Context, workoutID int) ([]workoutexercisejoin.WorkoutExercise, error) {
	args := m.Called(ctx, workoutID)
	return args.Get(0).([]workoutexercisejoin.WorkoutExercise), args.Error(1)
}
