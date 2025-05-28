package workout

import (
	"context"
	dto "workout-tracker/internal/dto/workout"
	join "workout-tracker/internal/model/workoutexercisejoin"
)

type FakeService struct {
	CreateErr   error
	UpdateErr   error
	DeleteErr   error
	AllErr      error
	AllResponse []dto.WorkoutWithExercises
	GetErr      error
	GetResponse *dto.WorkoutWithExercises
}

func (f *FakeService) CreateWorkout(ctx context.Context, userID int, name, title, category string,
	exercises []join.WorkoutExercise) error {
	return f.CreateErr
}
func (f *FakeService) UpdateWorkout(ctx context.Context, userID, workoutID int, name, title, category string,
	exercises []join.WorkoutExercise) error {
	return f.UpdateErr
}
func (f *FakeService) DeleteWorkout(ctx context.Context, userID, workoutID int) error {
	return f.DeleteErr
}
func (f *FakeService) GetAllWorkoutsWithExercises(ctx context.Context, userID int) ([]dto.WorkoutWithExercises, error) {
	return f.AllResponse, f.AllErr
}
func (f *FakeService) GetWorkoutByID(ctx context.Context, userID, workoutID int) (*dto.WorkoutWithExercises, error) {
	return f.GetResponse, f.GetErr
}
