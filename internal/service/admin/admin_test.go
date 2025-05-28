package admin

import (
	"errors"
	"testing"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/model/exercise"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func zapLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}

func TestCreateExercise_Success(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	req := dto.CreateExerciseRequest{Name: "Push-up", Description: "Upper body exercise"}
	mockRepo.On("CreateExercise", ctx, req).Return(42, nil)

	svc := NewAdminService(AdminServiceParams{
		Log:          zapLogger(),
		ExerciseRepo: mockRepo,
	})

	id, err := svc.CreateExercise(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 0, id)
	mockRepo.AssertExpectations(t)
}

func TestCreateExercise_Error(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	req := dto.CreateExerciseRequest{Name: "Squat"}
	errRepo := errors.New("db error")
	mockRepo.On("CreateExercise", ctx, req).Return(0, errRepo)

	svc := NewAdminService(AdminServiceParams{
		Log:          zapLogger(),
		ExerciseRepo: mockRepo,
	})

	id, err := svc.CreateExercise(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, 0, id)
	assert.Contains(t, err.Error(), "failed to create exercise in service")
	mockRepo.AssertExpectations(t)
}

func TestUpdateExercise_Success(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	input := dto.CreateExerciseRequest{Name: "Lunge"}
	mockRepo.On("UpdateExercise", ctx, 7, input).Return(nil)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})
	err := svc.UpdateExercise(ctx, 7, input)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateExercise_Error(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	input := dto.CreateExerciseRequest{Name: "Lunge"}
	errRepo := errors.New("not found")
	mockRepo.On("UpdateExercise", ctx, 7, input).Return(errRepo)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})
	err := svc.UpdateExercise(ctx, 7, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update exercise in service")
	mockRepo.AssertExpectations(t)
}

func TestGetAllExercises_Success(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	exs := []exercise.Exercise{{ID: 1, Name: "Plank"}, {ID: 2, Name: "Burpee"}}
	mockRepo.On("GetAllExercises", ctx).Return(exs, nil)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})

	result, err := svc.GetAllExercises(ctx)
	assert.NoError(t, err)
	assert.Equal(t, exs, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteExercise_Success(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	mockRepo.On("DeleteExercise", ctx, 3).Return(nil)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})
	err := svc.DeleteExercise(ctx, 3)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteExercise_Error(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	errRepo := errors.New("cannot delete")
	mockRepo.On("DeleteExercise", ctx, 3).Return(errRepo)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})
	err := svc.DeleteExercise(ctx, 3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete exercise")
	mockRepo.AssertExpectations(t)
}

func TestGetExerciseByID_Success(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	ex := &exercise.Exercise{ID: 9, Name: "Pull-up"}
	mockRepo.On("GetExerciseByID", ctx, 9).Return(ex, nil)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})

	result, err := svc.GetExerciseByID(ctx, 9)
	assert.NoError(t, err)
	assert.Equal(t, ex, result)
	mockRepo.AssertExpectations(t)
}

func TestGetExerciseByID_Error(t *testing.T) {
	ctx := t.Context()
	mockRepo := new(MockExerciseRepo)
	errRepo := errors.New("not found")
	mockRepo.On("GetExerciseByID", ctx, 9).Return((*exercise.Exercise)(nil), errRepo)

	svc := NewAdminService(AdminServiceParams{Log: zapLogger(), ExerciseRepo: mockRepo})

	result, err := svc.GetExerciseByID(ctx, 9)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get exercise")
	mockRepo.AssertExpectations(t)
}
