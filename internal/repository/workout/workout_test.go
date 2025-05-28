package workout

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

	model "workout-tracker/internal/model/workout"
	join "workout-tracker/internal/model/workoutexercisejoin"
)

var (
	pooler   *dockertest.Pool
	resource *dockertest.Resource
	pool     *pgxpool.Pool
	repo     *WorkoutRepository
)

func TestMain(m *testing.M) {
	t := &testing.T{} // временный T для Setenv
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")

	var err error
	pooler, err = dockertest.NewPool("")
	if err != nil {
		panic(err)
	}
	resource, err = pooler.Run("postgres", "13-alpine", []string{
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf(
		"postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		resource.GetPort("5432/tcp"),
		os.Getenv("DB_NAME"),
	)

	pooler.MaxWait = 30 * time.Second
	err = pooler.Retry(func() error {
		t.Setenv("DB_URL", dsn)
		dbObj, err := db.New(zap.NewNop().Sugar())
		if err != nil {
			return err
		}
		pool = dbObj.Pool
		return pool.Ping(context.Background())
	})
	if err != nil {
		panic(err)
	}

	exec := func(stmt string) {
		if _, err := pool.Exec(context.Background(), stmt); err != nil {
			panic(err)
		}
	}
	exec(`CREATE TABLE workouts (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		name TEXT,
		title TEXT,
		category TEXT,
		createdat TIMESTAMP,
		updatedat TIMESTAMP,
		deletedat TIMESTAMP
	);`)
	exec(`CREATE TABLE workout_exercise (
		workout_id INT,
		exercise_id INT,
		reps INT,
		sets INT
	);`)

	database, err := db.New(zap.NewNop().Sugar())
	if err != nil {
		panic(err)
	}

	repo = NewWorkoutRepository(WorkoutRepositoryParams{
		Log: zap.NewNop().Sugar(),
		DB:  database.Pool,
	})

	code := m.Run()

	pool.Close()
	_ = pooler.Purge(resource)
	os.Exit(code)
}

func TestCRUD_Workout(t *testing.T) {
	ctx := context.Background()
	now := time.Now().Truncate(time.Second)
	w := model.Workout{
		UserID:    42,
		Name:      "TestW",
		Title:     "TestTitle",
		Category:  "Cat",
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := repo.CreateWorkout(ctx, w)
	require.NoError(t, err)
	assert.Greater(t, id, 0)

	got, err := repo.GetWorkoutByID(ctx, id, w.UserID)
	require.NoError(t, err)
	assert.Equal(t, w.UserID, got.UserID)
	assert.Equal(t, w.Name, got.Name)

	got.Title = "NewTitle"
	got.UpdatedAt = now.Add(time.Minute)
	err = repo.UpdateWorkout(ctx, *got)
	require.NoError(t, err)

	updated, err := repo.GetWorkoutByID(ctx, id, w.UserID)
	require.NoError(t, err)
	assert.Equal(t, "NewTitle", updated.Title)

	list, err := repo.GetAllWorkouts(ctx, w.UserID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	err = repo.DeleteWorkout(ctx, id, w.UserID)
	require.NoError(t, err)

	_, err = repo.GetWorkoutByID(ctx, id, w.UserID)
	assert.Error(t, err)
}

func TestBulkAndFetchExercises(t *testing.T) {
	ctx := context.Background()
	id, err := repo.CreateWorkout(ctx, model.Workout{UserID: 1, Name: "W", Title: "T", Category: "C", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	require.NoError(t, err)

	exs := []join.WorkoutExercise{
		{WorkoutID: id, ExerciseID: 100, Reps: 10, Sets: 2},
		{WorkoutID: id, ExerciseID: 101, Reps: 5, Sets: 3},
	}

	err = repo.BulkInsertWorkoutExercises(ctx, exs)
	require.NoError(t, err)

	fetched, err := repo.GetWorkoutExercises(ctx, id)
	require.NoError(t, err)
	assert.Len(t, fetched, 2)
	assert.Equal(t, exs[0].ExerciseID, fetched[0].ExerciseID)

	err = repo.DeleteWorkoutExercises(ctx, id)
	require.NoError(t, err)

	emptied, err := repo.GetWorkoutExercises(ctx, id)
	require.NoError(t, err)
	assert.Len(t, emptied, 0)
}
