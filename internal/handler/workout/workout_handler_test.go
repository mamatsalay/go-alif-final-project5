package workout

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	dto "workout-tracker/internal/dto/workout"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupRouter(fs *FakeService) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", 7)
		c.Next()
	})
	h := NewWorkoutHandler(WorkoutHandlerParams{
		Service: fs,
		Logger:  zap.NewNop().Sugar(),
	})

	r.POST("/workouts", h.Create)
	r.PUT("/workouts/:id", h.Update)
	r.DELETE("/workouts/:id", h.Delete)
	r.GET("/workouts", h.GetAll)
	r.GET("/workouts/:id", h.Get)
	return r
}

func TestCreate_BadJSON(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/workouts", bytes.NewBufferString(`{invalid json}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	fs := &FakeService{CreateErr: errors.New("fail")}
	r := setupRouter(fs)
	payload := `{"name":"n","title":"t","category":"c","exercises":[{"exercise_id":1,"reps":5,"sets":2}]}`
	body := bytes.NewBufferString(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/workouts", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	fs := &FakeService{UpdateErr: errors.New("fail update")}
	r := setupRouter(fs)
	reqBody := dto.CreateWorkoutWithExercisesRequest{Name: "n"}
	bts, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/workouts/5", bytes.NewBuffer(bts))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	fs := &FakeService{UpdateErr: nil}
	r := setupRouter(fs)
	reqBody := dto.CreateWorkoutWithExercisesRequest{Name: "n"}
	bts, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PUT", "/workouts/5", bytes.NewBuffer(bts))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "workout updated", resp["message"])
}

func TestDelete_ServiceError(t *testing.T) {
	fs := &FakeService{DeleteErr: errors.New("fail delete")}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/workouts/5", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDelete_Success(t *testing.T) {
	fs := &FakeService{DeleteErr: nil}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/workouts/5", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "workout deleted", resp["message"])
}

func TestGetAll_Error(t *testing.T) {
	fs := &FakeService{AllErr: errors.New("fail all")}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/workouts", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet_InvalidID(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/workouts/abc", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGet_NotFound(t *testing.T) {
	fs := &FakeService{GetErr: errors.New("not found")}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/workouts/5", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
