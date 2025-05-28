package auth

import (
	"errors"
	"net/http"
	dto "workout-tracker/internal/dto/user"
	"workout-tracker/internal/erorrs"
	model "workout-tracker/internal/model/user"

	"go.uber.org/dig"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type AuthHandlerParams struct {
	dig.In

	Service AuthServiceInterface
	Logger  *zap.SugaredLogger
}

type AuthHandler struct {
	Service AuthServiceInterface
	Logger  *zap.SugaredLogger
}

func NewAuthHandler(p AuthHandlerParams) *AuthHandler {
	return &AuthHandler{
		Service: p.Service,
		Logger:  p.Logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest

	if err := c.Bind(&request); err != nil {
		h.Logger.Errorw("Invalid params: %s", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := h.Service.HashPassword(request.Password)
	if err != nil {
		h.Logger.Errorw("Error hashing password", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Username: request.Username,
		Password: hashed,
		Role:     model.UserRole,
	}
	id, err := h.Service.CreateUser(c.Request.Context(), user)
	if err != nil {
		h.Logger.Errorw("Error creating user", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	response := dto.RegisterResponse{
		ID:       id,
		Username: request.Username,
		Role:     model.UserRole,
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.Bind(&request); err != nil {
		h.Logger.Errorw("Invalid params: %s", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.GetUserByUsername(c.Request.Context(), request.Username)
	if err != nil {
		h.Logger.Errorw("Error getting user", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.CheckPassword(user.Password, request.Password); err != nil {
		h.Logger.Errorw("Error checking password", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, err := h.Service.GenerateAccessToken(user)
	if err != nil {
		h.Logger.Errorw("Error generating access token", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.Service.GenerateAndStoreRefreshToken(c.Request.Context(), user.ID)
	if err != nil {
		h.Logger.Errorw("Error generating refresh token", erorrs.ErrorKey, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token generation failed"})
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Username:     user.Username,
		Role:         user.Role,
		UserID:       user.ID,
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, newRefreshToken, err := h.Service.UpdateRefreshToken(c, req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, erorrs.ErrInternal) {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
