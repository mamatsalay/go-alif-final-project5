package workout

import (
	"context"
	model "workout-tracker/internal/model/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
)

type WorkoutRepositoryInterface interface {
	CreateWorkout(ctx context.Context, w model.Workout) (int, error)
	UpdateWorkout(ctx context.Context, w model.Workout) error
	DeleteWorkout(ctx context.Context, workoutID, userID int) error
	GetWorkoutByID(ctx context.Context, workoutID, userID int) (*model.Workout, error)
	BulkInsertWorkoutExercises(ctx context.Context, list []workoutexercisejoin.WorkoutExercise) error
	DeleteWorkoutExercises(ctx context.Context, workoutID int) error
	GetWorkoutExercises(ctx context.Context, workoutID int) ([]workoutexercisejoin.WorkoutExercise, error)
	GetAllWorkouts(ctx context.Context, userID int) ([]model.Workout, error)
}

var _ WorkoutRepositoryInterface = (*WorkoutRepository)(nil)
