package exercise

import (
	"context"
	"errors"
	"workout-tracker/internal/model/exercise"
	"workout-tracker/pkg/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type ExerciseRepositoryParams struct {
	dig.In
	DB  *db.DB
	Log *zap.SugaredLogger
}

type ExerciseRepository struct {
	Pool *pgxpool.Pool
	Log  *zap.SugaredLogger
}

func NewRepository(params ExerciseRepositoryParams) *ExerciseRepository {
	return &ExerciseRepository{
		Pool: params.DB.Pool,
		Log:  params.Log,
	}
}

// CreateExercise добавляет занятие в базу данных
func (r *ExerciseRepository) CreateExercise(ctx context.Context, ex exercise.Exercise) (exercise.Exercise, error) {
	if ex.Name == "" || ex.Sets == 0 || ex.Reps == 0 || ex.Weight == 0 || ex.Description == "" || ex.WorkoutID == 0 {
		return exercise.Exercise{}, errors.New("invalid parameters")
	}

	var created exercise.Exercise
	query := `
		INSERT INTO exercises (name, sets, reps, weight, description, workout_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, sets, reps, weight, description, workout_id, created_at, updated_at
	`
	err := r.Pool.QueryRow(ctx, query, ex.Name, ex.Sets, ex.Reps, ex.Weight, ex.Description, ex.WorkoutID).Scan(
		&created.ID, &created.Name, &created.Sets, &created.Reps, &created.Weight, &created.Description,
		&created.WorkoutID, &created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		r.Log.Errorw("Failed to create exercise", "error", err)
		return exercise.Exercise{}, err
	}
	return created, nil
}

// UpdateExercise обновляет занятие по ID
func (r *ExerciseRepository) UpdateExercise(ctx context.Context, ex exercise.Exercise) (exercise.Exercise, error) {
	if ex.Name == "" || ex.Sets == 0 || ex.Reps == 0 || ex.Weight == 0 || ex.Description == "" || ex.WorkoutID == 0 {
		return exercise.Exercise{}, errors.New("invalid parameters")
	}

	var updated exercise.Exercise
	query := `
		UPDATE exercises
		SET name=$1, sets=$2, reps=$3, weight=$4, description=$5, workout_id=$6, updated_at=NOW()
		WHERE id=$7
		RETURNING id, name, sets, reps, weight, description, workout_id, created_at, updated_at
	`
	err := r.Pool.QueryRow(ctx, query, ex.Name, ex.Sets, ex.Reps, ex.Weight, ex.Description, ex.WorkoutID, ex.ID).Scan(
		&updated.ID, &updated.Name, &updated.Sets, &updated.Reps, &updated.Weight, &updated.Description,
		&updated.WorkoutID, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		r.Log.Errorw("Failed to update exercise", "id", ex.ID, "error", err)
		return exercise.Exercise{}, err
	}
	return updated, nil
}

// GetAllExercises выводит все каталог всех занятий
func (r *ExerciseRepository) GetAllExercises(ctx context.Context) ([]exercise.Exercise, error) {
	query := `
		SELECT id, name, sets, reps, weight, description, workout_id, created_at, updated_at
		FROM exercises
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		r.Log.Errorw("Failed to query exercises", "error", err)
		return nil, err
	}
	defer rows.Close()

	var exercises []exercise.Exercise
	for rows.Next() {
		var ex exercise.Exercise
		err := rows.Scan(&ex.ID, &ex.Name, &ex.Sets, &ex.Reps, &ex.Weight, &ex.Description, &ex.WorkoutID, &ex.CreatedAt, &ex.UpdatedAt)
		if err != nil {
			r.Log.Errorw("Failed to scan exercise", "error", err)
			return nil, err
		}
		exercises = append(exercises, ex)
	}
	return exercises, nil
}
