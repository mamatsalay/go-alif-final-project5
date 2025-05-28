package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"workout-tracker/internal/erorrs"
	"workout-tracker/internal/handler"
	modeluser "workout-tracker/internal/model/user"
)

type FakeAuthService struct {
	User *modeluser.User
	Err  error
}

func (f *FakeAuthService) GetUserByUserID(ctx context.Context, id int) (*modeluser.User, error) {
	return f.User, f.Err
}

func setupRouter(secret string, fakeSvc *FakeAuthService) *gin.Engine {
	existing := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", existing)
	os.Setenv("JWT_SECRET", secret)

	m := handler.NewMiddleware(handler.MiddlewareParams{
		Log:     zap.NewNop().Sugar(),
		Service: fakeSvc,
	})

	r := gin.New()
	r.Use(m.AuthMiddleware())
	r.GET("/ok", func(c *gin.Context) {
		roleVal, _ := c.Get("role")
		c.JSON(http.StatusOK, gin.H{
			"userID": c.GetInt("userID"),
			"role":   roleVal,
		})
	})
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "missing or invalid token", body[erorrs.ErrorKey])
}

func TestAuthMiddleware_InvalidPrefix(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Token abc.def.ghi")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MissingClaims(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{})
	tok := jwt.New(jwt.SigningMethodHS256)
	tokStr, _ := tok.SignedString([]byte("secret"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer "+tokStr)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidClaimTypes(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{User: nil, Err: nil})
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "notanumber",
		"role":    123,
		"version": "noversion",
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)).Unix(),
	})
	tokStr, _ := tok.SignedString([]byte("secret"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer "+tokStr)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_UserNotFound(t *testing.T) {
	r := setupRouter("secret", &FakeAuthService{User: nil, Err: erorrs.ErrUserNotFound})
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
		"role":    string(modeluser.UserRole),
		"version": float64(0),
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)).Unix(),
	})
	tokStr, _ := tok.SignedString([]byte("secret"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer "+tokStr)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_VersionMismatch(t *testing.T) {
	userObj := &modeluser.User{ID: 1, Role: modeluser.UserRole, TokenVersion: 2}
	r := setupRouter("secret", &FakeAuthService{User: userObj, Err: nil})
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(1),
		"role":    string(userObj.Role),
		"version": float64(1),
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)).Unix(),
	})
	tokStr, _ := tok.SignedString([]byte("secret"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer "+tokStr)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Success(t *testing.T) {
	userObj := &modeluser.User{ID: 42, Role: modeluser.AdminRole, TokenVersion: 1}
	r := setupRouter("secret", &FakeAuthService{User: userObj, Err: nil})
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(42),
		"role":    string(modeluser.AdminRole),
		"version": float64(1),
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)).Unix(),
	})
	tokStr, _ := tok.SignedString([]byte("secret"))
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/ok", nil)
	req.Header.Set("Authorization", "Bearer "+tokStr)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(42), body["userID"])
	assert.Equal(t, string(modeluser.AdminRole), body["role"])
}

func TestAdminMiddleware_MissingRole(t *testing.T) {
	existing := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", existing)
	os.Setenv("JWT_SECRET", "secret")

	r := gin.New()
	r.Use(handler.NewMiddleware(handler.MiddlewareParams{Log: zap.NewNop().Sugar(), Service: &FakeAuthService{}}).AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminMiddleware_InvalidRoleType(t *testing.T) {
	existing := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", existing)
	os.Setenv("JWT_SECRET", "secret")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", 123) })
	r.Use(handler.NewMiddleware(handler.MiddlewareParams{Log: zap.NewNop().Sugar(), Service: &FakeAuthService{}}).AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminMiddleware_NotAdmin(t *testing.T) {
	existing := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", existing)
	os.Setenv("JWT_SECRET", "secret")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", modeluser.UserRole) })
	r.Use(handler.NewMiddleware(handler.MiddlewareParams{Log: zap.NewNop().Sugar(), Service: &FakeAuthService{}}).AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAdminMiddleware_Success(t *testing.T) {
	existing := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", existing)
	os.Setenv("JWT_SECRET", "secret")

	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("role", modeluser.AdminRole) })
	r.Use(handler.NewMiddleware(handler.MiddlewareParams{Log: zap.NewNop().Sugar(), Service: &FakeAuthService{}}).AdminMiddleware())
	r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
