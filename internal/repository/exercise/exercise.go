package exercise

import (
	"context"
	"errors"
	"time"

	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/erorrs"
	model "workout-tracker/internal/model/exercise"
	"workout-tracker/pkg/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type DBPool interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type ExerciseRepositoryParams struct {
	dig.In
	DB  *db.DB
	Log *zap.SugaredLogger
}

type ExerciseRepository struct {
	Pool DBPool
	Log  *zap.SugaredLogger
}

func NewRepository(params ExerciseRepositoryParams) *ExerciseRepository {
	return &ExerciseRepository{
		Pool: params.DB.Pool,
		Log:  params.Log,
	}
}

func (r *ExerciseRepository) CreateExercise(ctx context.Context, input dto.CreateExerciseRequest) (int, error) {
	var id int
	err := r.Pool.
		QueryRow(ctx,
			`INSERT INTO exercises(NAME, DESCRIPTION) VALUES ($1, $2) RETURNING id`,
			input.Name, input.Description).
		Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.Log.Errorw("exercise with this name already exists", "error", err)
			return 0, erorrs.ErrExerciseAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (r *ExerciseRepository) GetAllExercises(ctx context.Context) ([]model.Exercise, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id, name, description, createdat, updatedat
           FROM exercises
          WHERE deletedat IS NULL`,
	)
	if err != nil {
		r.Log.Errorw("error getting all exercises", "error", err)
		return nil, err
	}
	defer rows.Close()

	var result []model.Exercise
	for rows.Next() {
		var e model.Exercise
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt, &e.UpdatedAt); err != nil {
			r.Log.Errorw("error scanning exercise", "error", err)
			return nil, err
		}
		result = append(result, e)
	}
	if err := rows.Err(); err != nil {
		r.Log.Errorw("row iteration error", "error", err)
		return nil, err
	}
	return result, nil
}

func (r *ExerciseRepository) DeleteExercise(ctx context.Context, id int) error {
	_, err := r.Pool.Exec(ctx,
		"UPDATE exercises SET deletedat = $1 WHERE id = $2",
		time.Now(), id,
	)
	if err != nil {
		r.Log.Errorw("error deleting exercise", "id", id, "error", err)
	}
	return err
}

func (r *ExerciseRepository) GetExerciseByID(ctx context.Context, id int) (*model.Exercise, error) {
	var e model.Exercise
	err := r.Pool.
		QueryRow(ctx,
			`SELECT id, name, description, createdat, updatedat
               FROM exercises
              WHERE id = $1 AND deletedat IS NULL`,
			id).
		Scan(&e.ID, &e.Name, &e.Description, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		r.Log.Errorw("error getting exercise", "id", id, "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &e, nil
}

func (r *ExerciseRepository) UpdateExercise(ctx context.Context, id int, input dto.CreateExerciseRequest) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE exercises
            SET name = $1, description = $2, updatedat = $3
          WHERE id = $4`,
		input.Name, input.Description, time.Now(), id,
	)
	if err != nil {
		r.Log.Errorw("error updating exercise", "id", id, "error", err)
	}
	return err
}
