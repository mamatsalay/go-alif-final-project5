package workout

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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
	r.POST("/workouts/:id/photo", h.UpdatePhoto)
	r.GET("/workouts/:id/photo", h.GetPhoto)
	return r
}

func TestCreate_BadJSON(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts", bytes.NewBufferString(`{invalid json}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_ServiceError(t *testing.T) {
	fs := &FakeService{CreateErr: errors.New("fail")}
	r := setupRouter(fs)
	payload := `{"name":"n","title":"t","category":"c","exercises":[{"exercise_id":1,"reps":5,"sets":2}]}`
	body := bytes.NewBufferString(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts", body)
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
	req := httptest.NewRequest(http.MethodPut, "/workouts/5", bytes.NewBuffer(bts))
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
	req := httptest.NewRequest(http.MethodPut, "/workouts/5", bytes.NewBuffer(bts))
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
	req := httptest.NewRequest(http.MethodDelete, "/workouts/5", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDelete_Success(t *testing.T) {
	fs := &FakeService{DeleteErr: nil}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/workouts/5", http.NoBody)
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
	req := httptest.NewRequest(http.MethodGet, "/workouts", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGet_InvalidID(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/workouts/abc", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGet_NotFound(t *testing.T) {
	fs := &FakeService{GetErr: errors.New("not found")}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/workouts/5", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdatePhoto_InvalidID(t *testing.T) {
	fs := &FakeService{}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts/abc/photo", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePhoto_NoFile(t *testing.T) {
	fs := &FakeService{}
	r := setupRouter(fs)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts/5/photo", http.NoBody)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdatePhoto_SaveError(t *testing.T) {
	fs := &FakeService{UpdatePhotoErr: nil}
	r := setupRouter(fs)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("photo", "test.txt")
	fw.Write([]byte("data"))
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts/5/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdatePhoto_ServiceError(t *testing.T) {
	os.MkdirAll("uploads/workouts/5", 0755)
	defer os.RemoveAll("uploads")

	fs := &FakeService{UpdatePhotoErr: errors.New("svc fail")}
	r := setupRouter(fs)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("photo", "test.txt")
	fw.Write([]byte("data"))
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts/5/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdatePhoto_Success(t *testing.T) {
	os.MkdirAll("uploads/workouts/5", 0755)
	defer os.RemoveAll("uploads")

	fs := &FakeService{UpdatePhotoErr: nil}
	r := setupRouter(fs)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, _ := writer.CreateFormFile("photo", "test.txt")
	fw.Write([]byte("ok"))
	writer.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/workouts/5/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "photo uploaded successfully", resp["message"])
}

func TestGetPhoto_InvalidID(t *testing.T) {
	fs := &FakeService{}
	r := setupRouter(fs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/workouts/abc/photo", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPhoto_NotFoundInService(t *testing.T) {
	fs := &FakeService{GetErr: errors.New("not found")}
	r := setupRouter(fs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/workouts/5/photo", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
