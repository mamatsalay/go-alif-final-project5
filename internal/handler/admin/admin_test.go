package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"workout-tracker/internal/dto/exercise"
	exerciseRepsonse "workout-tracker/internal/model/exercise"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupRouter(svc *FakeAdminService) *gin.Engine {
	r := gin.Default()
	handler := NewAdminHandler(AdminHandlerParams{
		Service: svc,
		Logger:  zap.NewNop().Sugar(),
	})
	r.POST("/admin/exercises", handler.CreateExercise)
	r.PUT("/admin/exercises/:id", handler.UpdateExercise)
	r.GET("/admin/exercises", handler.GetAllExercises)
	r.DELETE("/admin/exercises/:id", handler.DeleteExercise)
	return r
}

func TestCreateExercise_BadJSON(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/exercises", bytes.NewBufferString("bad json"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateExercise_CreateError(t *testing.T) {
	r := setupRouter(&FakeAdminService{CreateErr: errors.New("db error")})
	body := exercise.CreateExerciseRequest{
		Name:        "Push Ups",
		Description: "Upper body",
	}
	data, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/exercises", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateExercise_Success(t *testing.T) {
	r := setupRouter(&FakeAdminService{CreateID: 42})
	body := exercise.CreateExerciseRequest{
		Name:        "Squat",
		Description: "Leg exercise",
	}
	data, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/admin/exercises", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]int
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 42, resp["exercise_id"])
}

func TestUpdateExercise_InvalidID(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/exercises/abc", bytes.NewBufferString(`{}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateExercise_BadJSON(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/exercises/1", bytes.NewBufferString("not json"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateExercise_Error(t *testing.T) {
	r := setupRouter(&FakeAdminService{UpdateErr: errors.New("update fail")})
	reqBody := exercise.CreateExerciseRequest{Name: "Deadlift", Description: "Back"}
	data, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/exercises/1", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateExercise_Success(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	reqBody := exercise.CreateExerciseRequest{Name: "Bench", Description: "Chest"}
	data, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/admin/exercises/2", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "exercise successfully updated", resp["answer"])
}

func TestGetAllExercises_Success(t *testing.T) {
	r := setupRouter(&FakeAdminService{
		GetAllResult: []exerciseRepsonse.Exercise{
			{ID: 1, Name: "Squat", Description: "Leg"},
			{ID: 2, Name: "Push-up", Description: "Upper body"},
		},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/exercises", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp []exerciseRepsonse.Exercise
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)
}

func TestGetAllExercises_Error(t *testing.T) {
	r := setupRouter(&FakeAdminService{GetAllErr: errors.New("db failure")})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/exercises", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteExercise_InvalidID(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/exercises/xyz", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteExercise_Error(t *testing.T) {
	r := setupRouter(&FakeAdminService{DeleteErr: errors.New("delete failed")})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/exercises/1", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteExercise_Success(t *testing.T) {
	r := setupRouter(&FakeAdminService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/admin/exercises/1", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)

	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "", resp["answer"])
}
