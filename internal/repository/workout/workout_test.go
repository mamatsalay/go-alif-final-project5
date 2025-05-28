package workout

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	model "workout-tracker/internal/model/workout"
	we "workout-tracker/internal/model/workoutexercisejoin"
)

func setupRepo(mockPool *MockPool) *WorkoutRepository {
	log := zaptest.NewLogger(&testing.T{}).Sugar()
	return &WorkoutRepository{Pool: mockPool, Log: log}
}

func TestCreateWorkout_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	ex := model.Workout{UserID: 1, Name: "n", Title: "t", Category: "c", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	row := new(MockRow)
	mp.On("QueryRow", ctx, mock.Anything,
		ex.UserID, ex.Name, ex.Title, ex.Category, ex.CreatedAt, ex.UpdatedAt).Return(row)
	row.On("Scan", mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
		*args.Get(0).(*int) = 55
	}).Return(nil)

	repo := setupRepo(mp)
	id, err := repo.CreateWorkout(ctx, ex)
	assert.NoError(t, err)
	assert.Equal(t, 55, id)
}

func TestCreateWorkout_Error(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	ex := model.Workout{UserID: 2, Name: "n", Title: "t", Category: "c", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	row := new(MockRow)
	mp.On("QueryRow", ctx, mock.Anything,
		ex.UserID, ex.Name, ex.Title, ex.Category, ex.CreatedAt, ex.UpdatedAt).Return(row)
	row.On("Scan", mock.Anything).Return(errors.New("fail"))

	repo := setupRepo(mp)
	id, err := repo.CreateWorkout(ctx, ex)
	assert.Error(t, err)
	assert.Zero(t, id)
}

func TestUpdateWorkout_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	w := model.Workout{ID: 10, UserID: 3, Name: "n", Title: "tt", Category: "cc", UpdatedAt: time.Now()}
	mp.On("Exec", ctx, mock.Anything,
		w.Title, w.Category, w.UpdatedAt, w.Name, w.ID, w.UserID).Return(pgconn.NewCommandTag(""), nil)

	repo := setupRepo(mp)
	err := repo.UpdateWorkout(ctx, w)
	assert.NoError(t, err)
}

func TestUpdateWorkout_Error(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	w := model.Workout{ID: 11, UserID: 4, Name: "n", Title: "tt", Category: "cc", UpdatedAt: time.Now()}
	mp.On("Exec", ctx, mock.Anything,
		w.Title, w.Category, w.UpdatedAt, w.Name, w.ID, w.UserID).Return(pgconn.NewCommandTag(""), errors.New("failup"))

	repo := setupRepo(mp)
	err := repo.UpdateWorkout(ctx, w)
	assert.Error(t, err)
}

func TestDeleteWorkout_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	mp.On("Exec", ctx, mock.Anything, mock.Anything, 5, 20).Return(pgconn.NewCommandTag(""), nil)

	repo := setupRepo(mp)
	err := repo.DeleteWorkout(ctx, 5, 20)
	assert.NoError(t, err)
}

func TestDeleteWorkout_Error(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	mp.On("Exec", ctx, mock.Anything, mock.Anything, 6, 21).Return(pgconn.NewCommandTag(""), errors.New("faildel"))

	repo := setupRepo(mp)
	err := repo.DeleteWorkout(ctx, 6, 21)
	assert.Error(t, err)
}

func TestGetWorkoutByID_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	w := model.Workout{ID: 7, UserID: 8, Name: "nm", Title: "tt", Category: "cc", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	row := new(MockRow)
	mp.On("QueryRow", ctx, mock.Anything, w.ID, w.UserID).Return(row)
	row.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	repo := setupRepo(mp)
	res, err := repo.GetWorkoutByID(ctx, w.ID, w.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestBulkInsertWorkoutExercises_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	list := []we.WorkoutExercise{{WorkoutID: 1, ExerciseID: 2, Reps: 3, Sets: 4}}
	mp.On("Exec", ctx, mock.Anything, list[0].WorkoutID, list[0].ExerciseID, list[0].Reps, list[0].Sets).
		Return(pgconn.NewCommandTag(""), nil)

	repo := setupRepo(mp)
	err := repo.BulkInsertWorkoutExercises(ctx, list)
	assert.NoError(t, err)
}

func TestBulkInsertWorkoutExercises_Error(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	list := []we.WorkoutExercise{{WorkoutID: 9, ExerciseID: 8, Reps: 7, Sets: 6}}
	mp.On("Exec", ctx, mock.Anything, 9, 8, 7, 6).Return(pgconn.NewCommandTag(""), errors.New("insfail"))

	repo := setupRepo(mp)
	err := repo.BulkInsertWorkoutExercises(ctx, list)
	assert.Error(t, err)
}

func TestDeleteWorkoutExercises_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	mp.On("Exec", ctx, mock.Anything, 42).Return(pgconn.NewCommandTag(""), nil)

	repo := setupRepo(mp)
	err := repo.DeleteWorkoutExercises(ctx, 42)
	assert.NoError(t, err)
}

func TestDeleteWorkoutExercises_Error(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	mp.On("Exec", ctx, mock.Anything, 43).Return(pgconn.NewCommandTag(""), errors.New("delfail"))

	repo := setupRepo(mp)
	err := repo.DeleteWorkoutExercises(ctx, 43)
	assert.Error(t, err)
}

func TestGetWorkoutExercises_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	r := new(MockRow)
	mp.On("Query", ctx, mock.Anything, 50).Return(r, nil)
	r.On("Next").Return(true).Once()
	r.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	r.On("Next").Return(false).Once()
	r.On("Close").Return()

	repo := setupRepo(mp)
	list, err := repo.GetWorkoutExercises(ctx, 50)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestGetWorkoutExercises_ScanError(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	r := new(MockRow)
	mp.On("Query", ctx, mock.Anything, 52).Return(r, nil)
	r.On("Next").Return(true)
	// scan for workout exercises expects 4 args
	r.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("scanerr"))
	r.On("Close").Return()

	repo := setupRepo(mp)
	list, err := repo.GetWorkoutExercises(ctx, 52)
	assert.Nil(t, list)
	assert.Error(t, err)
}

func TestGetAllWorkouts_Success(t *testing.T) {
	ctx := context.Background()
	mp := new(MockPool)
	r := new(MockRow)
	mp.On("Query", ctx, mock.Anything, 100).Return(r, nil)
	r.On("Next").Return(true).Once()
	r.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	r.On("Next").Return(false).Once()
	r.On("Close").Return()

	repo := setupRepo(mp)
	wos, err := repo.GetAllWorkouts(ctx, 100)
	assert.NoError(t, err)
	assert.Len(t, wos, 1)
}
