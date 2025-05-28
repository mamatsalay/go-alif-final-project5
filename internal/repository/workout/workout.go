package workout

import (
	"context"
	"fmt"
	"time"
	model "workout-tracker/internal/model/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
	"workout-tracker/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"go.uber.org/dig"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type WorkoutRepositoryParams struct {
	dig.In

	Log logger.SugaredLoggerInterface
	DB  DBPool
}

type WorkoutRepository struct {
	Log  logger.SugaredLoggerInterface
	Pool DBPool
}

func NewWorkoutRepository(params WorkoutRepositoryParams) WorkoutRepositoryInterface {
	return &WorkoutRepository{
		Log:  params.Log,
		Pool: params.DB,
	}
}

func (r *WorkoutRepository) CreateWorkout(ctx context.Context, input model.Workout) (int, error) {
	var id int
	err := r.Pool.QueryRow(ctx, `
		INSERT INTO workouts (user_id, name ,title, category, createdat, updatedat)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, input.UserID, input.Name, input.Title, input.Category, input.CreatedAt, input.UpdatedAt).Scan(&id)

	if err != nil {
		r.Log.Errorw("failed to create workout", "error", err)
		return 0, fmt.Errorf("create workout: %w", err)
	}
	return id, nil
}

func (r *WorkoutRepository) UpdateWorkout(ctx context.Context, workout model.Workout) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE workouts
		SET title = $1, category = $2, updatedat = $3, name = $4
		WHERE id = $5 AND user_id = $6 AND deletedat IS NULL
	`, workout.Title, workout.Category, workout.UpdatedAt, workout.Name, workout.ID, workout.UserID)

	if err != nil {
		r.Log.Errorw("failed to update workout", "error", err)
		return fmt.Errorf("update workout: %w", err)
	}
	return nil
}

func (r *WorkoutRepository) DeleteWorkout(ctx context.Context, workoutID int, userID int) error {
	now := time.Now()

	_, err := r.Pool.Exec(ctx, `
		UPDATE workouts
		SET deletedat = $1
		WHERE id = $2 AND user_id = $3 AND deletedat IS NULL
	`, now, workoutID, userID)

	if err != nil {
		r.Log.Errorw("failed to soft delete workout", "error", err)
		return fmt.Errorf("soft delete workout: %w", err)
	}
	return nil
}

func (r *WorkoutRepository) GetWorkoutByID(ctx context.Context, workoutID int, userID int) (*model.Workout, error) {
	var w model.Workout
	err := r.Pool.QueryRow(ctx, `
		SELECT id, user_id, name, title, category, createdat, updatedat
		FROM workouts
		WHERE id = $1 AND user_id = $2 AND deletedat IS NULL
	`, workoutID, userID).Scan(
		&w.ID,
		&w.UserID,
		&w.Name,
		&w.Title,
		&w.Category,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		r.Log.Errorw("failed to get workout", "error", err)
		return nil, fmt.Errorf("get workout: %w", err)
	}
	return &w, nil
}

func (r *WorkoutRepository) BulkInsertWorkoutExercises(ctx context.Context, list []workoutexercisejoin.WorkoutExercise) error {
	query := `
		INSERT INTO workout_exercise (workout_id, exercise_id, reps, sets)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range list {
		_, err := r.Pool.Exec(ctx, query,
			item.WorkoutID,
			item.ExerciseID,
			item.Reps,
			item.Sets,
		)
		if err != nil {
			r.Log.Errorw("failed to insert workout exercise", "error", err)
			return fmt.Errorf("insert workout exercise: %w", err)
		}
	}
	return nil
}

func (r *WorkoutRepository) DeleteWorkoutExercises(ctx context.Context, workoutID int) error {
	_, err := r.Pool.Exec(ctx, `
		DELETE FROM workout_exercise WHERE workout_id = $1
	`, workoutID)

	if err != nil {
		r.Log.Errorw("failed to delete workout exercises", "error", err)
		return fmt.Errorf("delete workout exercises: %w", err)
	}
	return nil
}

func (r *WorkoutRepository) GetWorkoutExercises(ctx context.Context, workoutID int) ([]workoutexercisejoin.WorkoutExercise, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT workout_id, exercise_id, reps, sets
		FROM workout_exercise
		WHERE workout_id = $1
	`, workoutID)

	if err != nil {
		r.Log.Errorw("failed to get workout exercises", "error", err)
		return nil, fmt.Errorf("get workout exercises: %w", err)
	}
	defer rows.Close()

	var list []workoutexercisejoin.WorkoutExercise
	for rows.Next() {
		var we workoutexercisejoin.WorkoutExercise
		if err := rows.Scan(&we.WorkoutID, &we.ExerciseID, &we.Reps, &we.Sets); err != nil {
			r.Log.Errorw("scan failed", "error", err)
			return nil, fmt.Errorf("get workout exercises: %w", err)
		}
		list = append(list, we)
	}

	return list, nil
}

func (r *WorkoutRepository) GetAllWorkouts(ctx context.Context, userID int) ([]model.Workout, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, user_id, name, title, category, createdat, updatedat
		FROM workouts
		WHERE user_id = $1 AND deletedat IS NULL
		ORDER BY createdat DESC
	`, userID)
	if err != nil {
		r.Log.Errorw("failed to fetch workouts", "error", err)
		return nil, fmt.Errorf("get workouts: %w", err)
	}
	defer rows.Close()

	var workouts []model.Workout
	for rows.Next() {
		var w model.Workout
		if err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.Title, &w.Category, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("get workouts: %w", err)
		}
		workouts = append(workouts, w)
	}

	return workouts, nil
}

func (r *WorkoutRepository) UpdateWorkoutPhoto(ctx context.Context, workoutID int, path string) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE workouts
		SET photo_path = $1, updatedat = NOW()
		WHERE id = $2
	`, path, workoutID)
	if err != nil {
		r.Log.Errorw("failed to update workout photo", "error", err)
		return fmt.Errorf("update workout photo: %w", err)
	}

	return nil
}
