package exercise

import (
	"context"
	dto "workout-tracker/internal/dto/exercise"
	model "workout-tracker/internal/model/exercise"
)

type ExerciseRepositoryInterface interface {
	CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error)
	GetAllExercises(ctx context.Context) ([]model.Exercise, error)
	DeleteExercise(ctx context.Context, id int) error
	GetExerciseByID(ctx context.Context, id int) (*model.Exercise, error)
	UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error
}

var _ ExerciseRepositoryInterface = (*ExerciseRepository)(nil)
