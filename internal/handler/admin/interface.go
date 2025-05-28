package admin

import (
	"context"
	"workout-tracker/internal/dto/exercise"
	exerciseRepsonse "workout-tracker/internal/model/exercise"
	"workout-tracker/internal/service/admin"
)

type AdminServiceInterface interface {
	CreateExercise(ctx context.Context, req exercise.CreateExerciseRequest) (int, error)
	UpdateExercise(ctx context.Context, id int, req exercise.CreateExerciseRequest) error
	GetAllExercises(ctx context.Context) ([]exerciseRepsonse.Exercise, error)
	DeleteExercise(ctx context.Context, id int) error
}

var _ AdminServiceInterface = (*admin.AdminService)(nil)
