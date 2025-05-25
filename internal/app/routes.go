package app

import (
	"workout-tracker/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *auth.AuthHandler) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
}
