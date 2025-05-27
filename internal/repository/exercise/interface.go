package exercise

import (
	"context"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/model/exercise"
)

type AdminServiceInterface interface {
	CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error)
	UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error
	GetAllExercises(ctx context.Context) ([]exercise.Exercise, error)
	DeleteExercise(ctx context.Context, id int) error
	GetExerciseByID(ctx context.Context, id int) (*exercise.Exercise, error)
}

var _ AdminServiceInterface = (*ExerciseRepository)(nil)
