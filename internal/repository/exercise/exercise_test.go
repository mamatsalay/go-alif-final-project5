package exercise

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	model "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/erorrs"
)

// Re-import mocks
// Assume mock.go is in same package path "workout-tracker/internal/repository/exercise"

func setupRepo(mockPool *MockPool) *ExerciseRepository {
	repo := &ExerciseRepository{Pool: mockPool, Log: zap.NewNop().Sugar()}
	return repo
}

func TestCreateExercise_Success(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `INSERT INTO exercises(NAME, DESCRIPTION) VALUES ($1, $2) RETURNING id`
	// stub QueryRow -> MockRow
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, "Name", "Desc").Return(mr)
	// Scan sets id
	mr.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
		*(args.Get(0).(*int)) = 42
	}).Return(nil)

	id, err := repo.CreateExercise(ctx, model.CreateExerciseRequest{Name: "Name", Description: "Desc"})
	assert.NoError(t, err)
	assert.Equal(t, 42, id)

	mp.AssertExpectations(t)
	mr.AssertExpectations(t)
}

func TestCreateExercise_Duplicate(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `INSERT INTO exercises(NAME, DESCRIPTION) VALUES ($1, $2) RETURNING id`
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, "N", "D").Return(mr)
	mr.On("Scan", mock.Anything).Return(&pgconn.PgError{Code: "23505"})

	id, err := repo.CreateExercise(ctx, model.CreateExerciseRequest{Name: "N", Description: "D"})
	assert.ErrorIs(t, err, erorrs.ErrExerciseAlreadyExists)
	assert.Zero(t, id)
}

func TestGetAllExercises(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
           FROM exercises
          WHERE deletedat IS NULL`
	mr := &MockRow{}
	mp.On("Query", ctx, sql).Return(mr, nil)
	// Next once true then false
	mr.On("Next").Return(true).Once()
	mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			*args.Get(0).(*int) = 1
			*args.Get(1).(*string) = "nm"
			*args.Get(2).(*string) = "ds"
			*args.Get(3).(*time.Time) = time.Now()
			*args.Get(4).(*time.Time) = time.Now()
		}).Return(nil).Once()
	mr.On("Next").Return(false).Once()
	mr.On("Err").Return(nil)
	mr.On("Close").Return()

	xs, err := repo.GetAllExercises(ctx)
	assert.NoError(t, err)
	assert.Len(t, xs, 1)
	assert.Equal(t, 1, xs[0].ID)
}

func TestDeleteExercise(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := "UPDATE exercises SET deletedat = $1 WHERE id = $2"
	mp.On("Exec", ctx, sql, mock.Anything, 5).Return(pgconn.NewCommandTag(""), nil)

	err := repo.DeleteExercise(ctx, 5)
	assert.NoError(t, err)
}

func TestDeleteExercise_Error(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := "UPDATE exercises SET deletedat = $1 WHERE id = $2"
	mp.On("Exec", ctx, sql, mock.Anything, 5).Return(pgconn.NewCommandTag(""), errors.New("fail"))

	err := repo.DeleteExercise(ctx, 5)
	assert.Error(t, err)
}

func TestGetExerciseByID(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
               FROM exercises
              WHERE id = $1 AND deletedat IS NULL`
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, 7).Return(mr)
	mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			*args.Get(0).(*int) = 7
			*args.Get(1).(*string) = "e"
			*args.Get(2).(*string) = "d"
			*args.Get(3).(*time.Time) = time.Now()
			*args.Get(4).(*time.Time) = time.Now()
		}).Return(nil)

	e, err := repo.GetExerciseByID(ctx, 7)
	assert.NoError(t, err)
	assert.Equal(t, 7, e.ID)
}

func TestGetExerciseByID_NotFound(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
               FROM exercises
              WHERE id = $1 AND deletedat IS NULL`
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, 8).Return(mr)
	// stub Scan to simulate not found
	mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(pgx.ErrNoRows)

	e, err := repo.GetExerciseByID(ctx, 8)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
	assert.Nil(t, e)
}

func TestUpdateExercise(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `UPDATE exercises
            SET name = $1, description = $2, updatedat = $3
          WHERE id = $4`
	mp.On("Exec", ctx, sql, "n", "d", mock.Anything, 9).Return(pgconn.NewCommandTag(""), nil)

	err := repo.UpdateExercise(ctx, 9, model.CreateExerciseRequest{Name: "n", Description: "d"})
	assert.NoError(t, err)
}

func TestUpdateExercise_Error(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `UPDATE exercises
            SET name = $1, description = $2, updatedat = $3
          WHERE id = $4`
	mp.On("Exec", ctx, sql, "n", "d", mock.Anything, 9).Return(pgconn.NewCommandTag(""), errors.New("err"))

	err := repo.UpdateExercise(ctx, 9, model.CreateExerciseRequest{Name: "n", Description: "d"})
	assert.Error(t, err)
}

func TestCreateExercise_OtherError(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `INSERT INTO exercises(NAME, DESCRIPTION) VALUES ($1, $2) RETURNING id`
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, "X", "Y").Return(mr)
	mr.On("Scan", mock.Anything).Return(errors.New("oops"))

	id, err := repo.CreateExercise(ctx, model.CreateExerciseRequest{Name: "X", Description: "Y"})
	assert.EqualError(t, err, "oops")
	assert.Zero(t, id)
}

func TestGetAllExercises_QueryError(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
           FROM exercises
          WHERE deletedat IS NULL`
	var rows pgx.Rows = (*MockRow)(nil) // typed nil to avoid panic
	mp.On("Query", ctx, sql).Return(rows, errors.New("qerr"))

	xs, err := repo.GetAllExercises(ctx)
	assert.Nil(t, xs)
	assert.EqualError(t, err, "qerr")
}

func TestGetAllExercises_ScanError(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
           FROM exercises
          WHERE deletedat IS NULL`
	mr := &MockRow{}
	mp.On("Query", ctx, sql).Return(mr, nil)
	mr.On("Next").Return(true)
	mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("scanfail"))
	mr.On("Close").Return()

	xs, err := repo.GetAllExercises(ctx)
	assert.Nil(t, xs)
	assert.EqualError(t, err, "scanfail")
}

func TestGetAllExercises_RowsErr(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
           FROM exercises
          WHERE deletedat IS NULL`
	mr := &MockRow{}
	mp.On("Query", ctx, sql).Return(mr, nil)
	mr.On("Next").Return(false)
	mr.On("Err").Return(errors.New("itererr"))
	mr.On("Close").Return()

	xs, err := repo.GetAllExercises(ctx)
	assert.Nil(t, xs)
	assert.EqualError(t, err, "itererr")
}

func TestGetExerciseByID_ScanError(t *testing.T) {
	mp := &MockPool{}
	repo := setupRepo(mp)

	ctx := context.Background()
	sql := `SELECT id, name, description, createdat, updatedat
               FROM exercises
              WHERE id = $1 AND deletedat IS NULL`
	mr := &MockRow{}
	mp.On("QueryRow", ctx, sql, 10).Return(mr)
	mr.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("scanidfail"))

	e, err := repo.GetExerciseByID(ctx, 10)
	assert.Nil(t, e)
	assert.EqualError(t, err, "scanidfail")
}
