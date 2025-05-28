package workout

import (
	"context"
	dto "workout-tracker/internal/dto/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
	"workout-tracker/internal/service/workout"
)

type WorkoutServiceInterface interface {
	CreateWorkout(ctx context.Context, userID int, name, title, category string, exercises []workoutexercisejoin.WorkoutExercise) error
	UpdateWorkout(ctx context.Context, userID, workoutID int, name, title, category string, exercises []workoutexercisejoin.WorkoutExercise) error
	DeleteWorkout(ctx context.Context, userID, workoutID int) error
	GetAllWorkoutsWithExercises(ctx context.Context, userID int) ([]dto.WorkoutWithExercises, error)
	GetWorkoutByID(ctx context.Context, userID, workoutID int) (*dto.WorkoutWithExercises, error)
	UpdateWorkoutPhoto(ctx context.Context, workoutID int, path string) error
}

var _ WorkoutServiceInterface = (*workout.WorkoutService)(nil)
