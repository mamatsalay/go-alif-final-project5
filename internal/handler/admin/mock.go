package admin

import (
	"context"

	"workout-tracker/internal/dto/exercise"
	exerciseRepsonse "workout-tracker/internal/model/exercise"
)

type FakeAdminService struct {
	CreateID     int
	CreateErr    error
	UpdateErr    error
	GetAllResult []exerciseRepsonse.Exercise
	GetAllErr    error
	DeleteErr    error
}

func (f *FakeAdminService) CreateExercise(ctx context.Context, req exercise.CreateExerciseRequest) (int, error) {
	return f.CreateID, f.CreateErr
}

func (f *FakeAdminService) UpdateExercise(ctx context.Context, id int, req exercise.CreateExerciseRequest) error {
	return f.UpdateErr
}

func (f *FakeAdminService) GetAllExercises(ctx context.Context) ([]exerciseRepsonse.Exercise, error) {
	return f.GetAllResult, f.GetAllErr
}

func (f *FakeAdminService) DeleteExercise(ctx context.Context, id int) error {
	return f.DeleteErr
}
