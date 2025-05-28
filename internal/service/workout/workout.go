package workout

import (
	"context"
	"fmt"
	"time"
	dto "workout-tracker/internal/dto/workout"
	model "workout-tracker/internal/model/workout"
	joinModel "workout-tracker/internal/model/workoutexercisejoin"
	workoutInterface "workout-tracker/internal/repository/workout"
	"workout-tracker/pkg/logger"

	"go.uber.org/dig"
)

type WorkoutServiceParams struct {
	dig.In

	Repo workoutInterface.WorkoutRepositoryInterface
	Log  logger.SugaredLoggerInterface
}

type WorkoutService struct {
	Repo workoutInterface.WorkoutRepositoryInterface
	Log  logger.SugaredLoggerInterface
}

func NewWorkoutService(params WorkoutServiceParams) *WorkoutService {
	return &WorkoutService{
		Repo: params.Repo,
		Log:  params.Log,
	}
}

func (s *WorkoutService) CreateWorkout(ctx context.Context, userID int, name, title, category string, exercises []joinModel.WorkoutExercise) error {
	now := time.Now()
	workout := model.Workout{
		UserID:    userID,
		Name:      name,
		Title:     title,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.Repo.CreateWorkout(ctx, workout)
	if err != nil {
		s.Log.Errorw("failed to create workout", "error", err)
		return fmt.Errorf("create workout: %w", err)
	}

	for i := range exercises {
		exercises[i].WorkoutID = id
	}

	err = s.Repo.BulkInsertWorkoutExercises(ctx, exercises)
	if err != nil {
		s.Log.Errorw("failed to insert exercises", "error", err)
		return fmt.Errorf("insert exercises: %w", err)
	}

	return nil
}

func (s *WorkoutService) UpdateWorkout(ctx context.Context, userID, workoutID int,
	name, title, category string, exercises []joinModel.WorkoutExercise) error {
	workout := model.Workout{
		ID:        workoutID,
		Name:      name,
		UserID:    userID,
		Title:     title,
		Category:  category,
		UpdatedAt: time.Now(),
	}

	if err := s.Repo.UpdateWorkout(ctx, workout); err != nil {
		return fmt.Errorf("update workout: %w", err)
	}

	if err := s.Repo.DeleteWorkoutExercises(ctx, workoutID); err != nil {
		return fmt.Errorf("delete workout exercises: %w", err)
	}

	for i := range exercises {
		exercises[i].WorkoutID = workoutID
	}

	err := s.Repo.BulkInsertWorkoutExercises(ctx, exercises)
	if err != nil {
		s.Log.Errorw("failed to insert exercises", "error", err)
		return fmt.Errorf("insert exercises: %w", err)
	}

	return nil
}

func (s *WorkoutService) DeleteWorkout(ctx context.Context, userID, workoutID int) error {
	err := s.Repo.DeleteWorkout(ctx, workoutID, userID)
	if err != nil {
		s.Log.Errorw("failed to delete workout", "error", err)
		return fmt.Errorf("delete workout: %w", err)
	}

	return nil
}

func (s *WorkoutService) GetAllWorkoutsWithExercises(ctx context.Context, userID int) ([]dto.WorkoutWithExercises, error) {
	workouts, err := s.Repo.GetAllWorkouts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get all workouts: %w", err)
	}

	var result []dto.WorkoutWithExercises
	for _, w := range workouts {
		exercises, err := s.Repo.GetWorkoutExercises(ctx, w.ID)
		if err != nil {
			s.Log.Errorw("failed to fetch exercises for workout", "workout_id", w.ID, "error", err)
			return nil, fmt.Errorf("fetch exercises for workout: %w", err)
		}

		result = append(result, dto.WorkoutWithExercises{
			Workout:   w,
			Exercises: exercises,
		})
	}

	return result, nil
}

func (s *WorkoutService) GetWorkoutByID(ctx context.Context, userID int, workoutID int) (*dto.WorkoutWithExercises, error) {
	workout, err := s.Repo.GetWorkoutByID(ctx, workoutID, userID)
	if err != nil {
		s.Log.Errorw("failed to get workout", "workoutID", workoutID, "error", err)
		return nil, fmt.Errorf("get workout: %w", err)
	}

	exercises, err := s.Repo.GetWorkoutExercises(ctx, workoutID)
	if err != nil {
		s.Log.Errorw("failed to get exercises for workout", "workoutID", workoutID, "error", err)
		return nil, fmt.Errorf("get exercises for workout: %w", err)
	}

	return &dto.WorkoutWithExercises{
		Workout:   *workout,
		Exercises: exercises,
	}, nil
}

func (s *WorkoutService) UpdateWorkoutPhoto(ctx context.Context, workoutID int, path string) error {
	err := s.Repo.UpdateWorkoutPhoto(ctx, workoutID, path)
	if err != nil {
		return fmt.Errorf("update workout photo: %w", err)
	}

	return nil
}
