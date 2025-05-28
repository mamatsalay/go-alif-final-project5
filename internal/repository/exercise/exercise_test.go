package exercise

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/erorrs"

	"workout-tracker/pkg/db"

	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var exerciseRepo *ExerciseRepository

func TestMain(m *testing.M) {
	t := &testing.T{} // initialize for Setenv
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")

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
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		resource.GetPort("5432/tcp"), os.Getenv("DB_NAME"),
	)

	pooler.MaxWait = 30 * time.Second
	err = pooler.Retry(func() error {
		t.Setenv("DB_URL", dsn)
		database, err := db.New(zap.NewNop().Sugar())
		if err != nil {
			return err
		}
		return database.Pool.Ping(context.Background())
	})
	if err != nil {
		panic(err)
	}

	database, err := db.New(zap.NewNop().Sugar())
	if err != nil {
		panic(err)
	}

	exec := func(q string) {
		if _, err := database.Pool.Exec(context.Background(), q); err != nil {
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
		DB:  database,
		Log: zap.NewNop().Sugar(),
	})

	code := m.Run()

	database.Pool.Close()
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

func TestCreateExercise_DuplicateName(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Sit-up", Description: "Core"}
	_, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)

	_, err = exerciseRepo.CreateExercise(ctx, req)
	assert.ErrorIs(t, err, erorrs.ErrExerciseAlreadyExists)
}

func TestGetNonexistentExercise(t *testing.T) {
	ctx := context.Background()
	_, err := exerciseRepo.GetExerciseByID(ctx, 9999)
	assert.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestDeleteExercise(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Plank", Description: "Core hold"}
	id, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)

	err = exerciseRepo.DeleteExercise(ctx, id)
	require.NoError(t, err)

	_, err = exerciseRepo.GetExerciseByID(ctx, id)
	assert.Error(t, err)
}

func TestUpdateExercise(t *testing.T) {
	ctx := context.Background()
	req := dto.CreateExerciseRequest{Name: "Lunge", Description: "Legs"}
	id, err := exerciseRepo.CreateExercise(ctx, req)
	require.NoError(t, err)

	update := dto.CreateExerciseRequest{
		Name:        "Lunge Updated",
		Description: "Leg strength",
	}
	err = exerciseRepo.UpdateExercise(ctx, id, update)
	require.NoError(t, err)

	e, err := exerciseRepo.GetExerciseByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, update.Name, e.Name)
	assert.Equal(t, update.Description, e.Description)
}

func TestGetAllExercises(t *testing.T) {
	ctx := context.Background()
	exercises := []dto.CreateExerciseRequest{
		{Name: "Deadlift", Description: "Back"},
		{Name: "Bench Press", Description: "Chest"},
	}
	for _, ex := range exercises {
		_, err := exerciseRepo.CreateExercise(ctx, ex)
		require.NoError(t, err)
	}

	all, err := exerciseRepo.GetAllExercises(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(all), 2)
	for _, ex := range exercises {
		found := false
		for _, got := range all {
			if got.Name == ex.Name && got.Description == ex.Description {
				found = true
				break
			}
		}
		assert.True(t, found, "exercise not found: %s", ex.Name)
	}
}
