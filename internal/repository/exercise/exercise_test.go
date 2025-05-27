package exercise

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"workout-tracker/pkg/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/erorrs"
)

var pool *pgxpool.Pool
var exerciseRepo *ExerciseRepository

func TestMain(m *testing.M) {
	pooler, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	resource, err := pooler.Run("postgres", "13-alpine", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("postgres://postgres:secret@localhost:%s/testdb?sslmode=disable", resource.GetPort("5432/tcp"))
	pooler.MaxWait = 30 * time.Second
	err = pooler.Retry(func() error {
		var err error
		pool, err = pgxpool.New(context.Background(), dsn)
		if err != nil {
			return err
		}
		return pool.Ping(context.Background())
	})
	if err != nil {
		panic(err)
	}

	exec := func(q string) {
		if _, err := pool.Exec(context.Background(), q); err != nil {
			panic(err)
		}
	}
	exec(`CREATE TABLE exercises (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT,
	createdat TIMESTAMP DEFAULT NOW(),
	updatedat TIMESTAMP DEFAULT NOW(),
	deletedat TIMESTAMP
);`)

	exerciseRepo = NewRepository(ExerciseRepositoryParams{
		DB:  (*db.DB)(&struct{ Pool *pgxpool.Pool }{Pool: pool}),
		Log: zap.NewNop().Sugar(),
	})

	code := m.Run()

	pool.Close()
	_ = pooler.Purge(resource)
	os.Exit(code)
}

func TestCreateAndGetExercise(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Push-up", Description: "Upper body"}
	id, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	e, err := exerciseRepo.GetExerciseByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, id, e.ID)
	assert.Equal(t, req.Name, e.Name)
	assert.Equal(t, req.Description, e.Description)
}

func TestCreateExercise_Duplicate(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Squat", Description: "Legs"}
	_, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)
	_, err = exerciseRepo.CreateExercise(ctx, req)
	assert.ErrorIs(t, err, erorrs.ErrExerciseAlreadyExists)
}

func TestGetAllExercises(t *testing.T) {
	ctx := context.Background()
	exerciseRepo.CreateExercise(ctx, dto.CreateExerciseRequest{Name: "A1"})
	exerciseRepo.CreateExercise(ctx, dto.CreateExerciseRequest{Name: "A2"})

	list, err := exerciseRepo.GetAllExercises(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 2)
	for _, e := range list {
		assert.NotZero(t, e.ID)
		assert.NotEmpty(t, e.Name)
	}
}

func TestUpdateAndDeleteExercise(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Lunge", Description: "Legs"}
	id, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)

	update := dto.CreateExerciseRequest{Name: "Lunge modified", Description: "Leg muscles"}
	err = exerciseRepo.UpdateExercise(ctx, id, update)
	require.NoError(t, err)
	e, err := exerciseRepo.GetExerciseByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, update.Name, e.Name)

	err = exerciseRepo.DeleteExercise(ctx, id)
	require.NoError(t, err)
	_, err = exerciseRepo.GetExerciseByID(ctx, id)
	assert.Error(t, err)
}
