package workout

import (
	"context"
	model "workout-tracker/internal/model/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
)

type WorkoutRepositoryInterface interface {
	CreateWorkout(ctx context.Context, input model.Workout) (int, error)
	UpdateWorkout(ctx context.Context, workout model.Workout) error
	DeleteWorkout(ctx context.Context, workoutID int, userID int) error
	GetWorkoutByID(ctx context.Context, workoutID int, userID int) (*model.Workout, error)
	BulkInsertWorkoutExercises(ctx context.Context, list []workoutexercisejoin.WorkoutExercise) error
	DeleteWorkoutExercises(ctx context.Context, workoutID int) error
	GetWorkoutExercises(ctx context.Context, workoutID int) ([]workoutexercisejoin.WorkoutExercise, error)
	GetAllWorkouts(ctx context.Context, userID int) ([]model.Workout, error)
}

var _ WorkoutRepositoryInterface = (*WorkoutRepository)(nil)
