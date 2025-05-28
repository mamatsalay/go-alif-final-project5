package workout

import (
	"context"
	"fmt"

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
	return fmt.Errorf("error bulk insert workout: %w", args.Error(0))
}

func (m *WorkoutRepoMock) UpdateWorkout(ctx context.Context, w workout.Workout) error {
	args := m.Called(ctx, w)
	return fmt.Errorf("error: update workout%w", args.Error(0))
}

func (m *WorkoutRepoMock) DeleteWorkoutExercises(ctx context.Context, workoutID int) error {
	args := m.Called(ctx, workoutID)
	return fmt.Errorf("error: delete workout ex%w", args.Error(0))
}

func (m *WorkoutRepoMock) DeleteWorkout(ctx context.Context, workoutID, userID int) error {
	args := m.Called(ctx, workoutID, userID)
	return fmt.Errorf("error: delete workout%w", args.Error(0))
}

func (m *WorkoutRepoMock) GetAllWorkouts(ctx context.Context, userID int) ([]workout.Workout, error) {
	args := m.Called(ctx, userID)

	workouts, ok := args.Get(0).([]workout.Workout)
	if !ok {
		return nil, fmt.Errorf("invalid type for []workout.Workout: %w", args.Error(0))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error in GetAllWorkouts: %w", err)
	}

	return workouts, nil
}

func (m *WorkoutRepoMock) GetWorkoutByID(ctx context.Context, workoutID, userID int) (*workout.Workout, error) {
	args := m.Called(ctx, workoutID, userID)

	workoutObj, ok := args.Get(0).(*workout.Workout)
	if !ok {
		return nil, fmt.Errorf("invalid type for *workout.Workout: %w", args.Error(0))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error in GetWorkoutByID: %w", err)
	}

	return workoutObj, nil
}

func (m *WorkoutRepoMock) GetWorkoutExercises(ctx context.Context, workoutID int) ([]workoutexercisejoin.WorkoutExercise, error) {
	args := m.Called(ctx, workoutID)

	exercises, ok := args.Get(0).([]workoutexercisejoin.WorkoutExercise)
	if !ok {
		return nil, fmt.Errorf("invalid type for []workoutexercisejoin.WorkoutExercise: %w", args.Error(0))
	}

	if err := args.Error(1); err != nil {
		return nil, fmt.Errorf("error in GetWorkoutExercises: %w", err)
	}

	return exercises, nil
}
