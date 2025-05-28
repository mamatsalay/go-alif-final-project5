package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	dto "workout-tracker/internal/dto/user"
	model "workout-tracker/internal/model/user"
)

func setupRouter(fs AuthServiceInterface) *gin.Engine {
	r := gin.New()
	h := NewAuthHandler(AuthHandlerParams{Service: fs, Logger: zap.NewNop().Sugar()})

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/refresh", h.RefreshToken)
	return r
}

func TestRegister_BadBind(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("notjson"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_HashError(t *testing.T) {
	fs := &FakeService{HashErr: errors.New("hash fail")}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_CreateError(t *testing.T) {
	fs := &FakeService{CreatedID: 0, CreateErr: errors.New("dup")}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_Success(t *testing.T) {
	fs := &FakeService{CreatedID: 5}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data dto.RegisterResponse `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 5, resp.Data.ID)
}

func TestLogin_BadBind(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString("x"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_GetUserError(t *testing.T) {
	fs := &FakeService{FindErr: errors.New("no user")}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_CheckPasswordError(t *testing.T) {
	fs := &FakeService{FoundUser: &model.User{Username: "u", Password: "h"}, PasswordCheckErr: errors.New("wrong")}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_TokenError(t *testing.T) {
	fs := &FakeService{FoundUser: &model.User{ID: 1, Username: "u", Password: "h"}, AccessErr: errors.New("tokfail")}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	fs := &FakeService{FoundUser: &model.User{ID: 2, Username: "u", Password: "h", Role: model.UserRole}, AccessToken: "at", RefreshToken: "rt"}
	r := setupRouter(fs)
	payload := `{"username":"u","password":"p"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp dto.LoginResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "", resp.AccessToken)
}

func TestRefresh_BadJSON(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString("{}"))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefresh_Error(t *testing.T) {
	fs := &FakeService{UpdateErr: errors.New("bad")}
	r := setupRouter(fs)
	payload := `{"refresh_token":"rt"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefresh_Success(t *testing.T) {
	fs := &FakeService{UpdateErr: nil}
	r := setupRouter(fs)
	payload := `{"refresh_token":"rt"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(payload))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "newAccess", resp["access_token"])
	assert.Equal(t, "newRefresh", resp["refresh_token"])
}

func TestLogin_RefreshGenerationError(t *testing.T) {
	fs := &FakeService{FoundUser: &model.User{ID: 10, Username: "bob", Password: "h", Role: model.UserRole}, AccessToken: "at", RefreshErr: errors.New("storefail")}
	r := setupRouter(fs)
	payload := `{"username":"bob","password":"pw"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRegister_BadContentType(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString(`{"username":"u","password":"p"}`))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRefresh_MissingField(t *testing.T) {
	r := setupRouter(&FakeService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(`{"token":"rt"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
