package workout_test

import (
	"context"
	"errors"
	"testing"

	model "workout-tracker/internal/model/workout"
	joinModel "workout-tracker/internal/model/workoutexercisejoin"
	"workout-tracker/internal/service/workout"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

type stubRepo struct {
	CreateWorkoutFn              func(ctx context.Context, w model.Workout) (int, error)
	BulkInsertWorkoutExercisesFn func(ctx context.Context, ex []joinModel.WorkoutExercise) error
	UpdateWorkoutFn              func(ctx context.Context, w model.Workout) error
	DeleteWorkoutExercisesFn     func(ctx context.Context, workoutID int) error
	DeleteWorkoutFn              func(ctx context.Context, workoutID int, userID int) error
	GetAllWorkoutsFn             func(ctx context.Context, userID int) ([]model.Workout, error)
	GetWorkoutExercisesFn        func(ctx context.Context, workoutID int) ([]joinModel.WorkoutExercise, error)
	GetWorkoutByIDFn             func(ctx context.Context, workoutID int, userID int) (*model.Workout, error)
	UpdateWorkoutPhotoFn         func(ctx context.Context, workoutID int, path string) error
}

func (s *stubRepo) UpdateWorkoutPhoto(ctx context.Context, workoutID int, path string) error {
	return s.UpdateWorkoutPhotoFn(ctx, workoutID, path)
}

func (s *stubRepo) CreateWorkout(ctx context.Context, w model.Workout) (int, error) {
	return s.CreateWorkoutFn(ctx, w)
}
func (s *stubRepo) BulkInsertWorkoutExercises(ctx context.Context, ex []joinModel.WorkoutExercise) error {
	return s.BulkInsertWorkoutExercisesFn(ctx, ex)
}
func (s *stubRepo) UpdateWorkout(ctx context.Context, w model.Workout) error {
	return s.UpdateWorkoutFn(ctx, w)
}
func (s *stubRepo) DeleteWorkoutExercises(ctx context.Context, workoutID int) error {
	return s.DeleteWorkoutExercisesFn(ctx, workoutID)
}
func (s *stubRepo) DeleteWorkout(ctx context.Context, workoutID int, userID int) error {
	return s.DeleteWorkoutFn(ctx, workoutID, userID)
}
func (s *stubRepo) GetAllWorkouts(ctx context.Context, userID int) ([]model.Workout, error) {
	return s.GetAllWorkoutsFn(ctx, userID)
}
func (s *stubRepo) GetWorkoutExercises(ctx context.Context, workoutID int) ([]joinModel.WorkoutExercise, error) {
	return s.GetWorkoutExercisesFn(ctx, workoutID)
}
func (s *stubRepo) GetWorkoutByID(ctx context.Context, workoutID int, userID int) (*model.Workout, error) {
	return s.GetWorkoutByIDFn(ctx, workoutID, userID)
}

func newTestService(t *testing.T, repo *stubRepo) *workout.WorkoutService {
	t.Helper()
	logger := zaptest.NewLogger(t).Sugar()
	return workout.NewWorkoutService(workout.WorkoutServiceParams{Repo: repo, Log: logger})
}

func TestCreateWorkout_Success(t *testing.T) {
	repo := &stubRepo{
		CreateWorkoutFn: func(ctx context.Context, w model.Workout) (int, error) {
			return 1, nil
		},
		BulkInsertWorkoutExercisesFn: func(ctx context.Context, ex []joinModel.WorkoutExercise) error {
			return nil
		},
	}
	service := newTestService(t, repo)
	exercises := []joinModel.WorkoutExercise{{ExerciseID: 1, Sets: 3}}
	err := service.CreateWorkout(t.Context(), 1, "Test", "Title", "Strength", exercises)
	assert.NoError(t, err)
}

func TestCreateWorkout_BulkInsertFails(t *testing.T) {
	repo := &stubRepo{
		CreateWorkoutFn: func(ctx context.Context, w model.Workout) (int, error) {
			return 1, nil
		},
		BulkInsertWorkoutExercisesFn: func(ctx context.Context, ex []joinModel.WorkoutExercise) error {
			return errors.New("bulk insert error")
		},
	}
	service := newTestService(t, repo)
	err := service.CreateWorkout(t.Context(), 1, "Test", "Title", "Strength", nil)
	assert.Error(t, err)
}

func TestDeleteWorkout_Success(t *testing.T) {
	repo := &stubRepo{
		DeleteWorkoutFn: func(ctx context.Context, workoutID int, userID int) error {
			return nil
		},
	}
	service := newTestService(t, repo)
	err := service.DeleteWorkout(t.Context(), 1, 1)
	assert.NoError(t, err)
}

func TestGetWorkoutByID_Success(t *testing.T) {
	repo := &stubRepo{
		GetWorkoutByIDFn: func(ctx context.Context, workoutID int, userID int) (*model.Workout, error) {
			return &model.Workout{ID: 1, Name: "Test"}, nil
		},
		GetWorkoutExercisesFn: func(ctx context.Context, workoutID int) ([]joinModel.WorkoutExercise, error) {
			return []joinModel.WorkoutExercise{{ExerciseID: 1}}, nil
		},
	}
	service := newTestService(t, repo)
	res, err := service.GetWorkoutByID(t.Context(), 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Workout.ID)
}

func TestGetAllWorkoutsWithExercises_FetchFails(t *testing.T) {
	repo := &stubRepo{
		GetAllWorkoutsFn: func(ctx context.Context, userID int) ([]model.Workout, error) {
			return []model.Workout{{ID: 1}}, nil
		},
		GetWorkoutExercisesFn: func(ctx context.Context, workoutID int) ([]joinModel.WorkoutExercise, error) {
			return nil, errors.New("fetch error")
		},
	}
	service := newTestService(t, repo)
	_, err := service.GetAllWorkoutsWithExercises(t.Context(), 1)
	assert.Error(t, err)
}
