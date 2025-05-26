package app

import (
	"workout-tracker/internal/handler"
	"workout-tracker/internal/handler/admin"
	"workout-tracker/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *auth.AuthHandler, a *admin.AdminHandler, m *handler.Middleware) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	admin := r.Group("/admin").Use(m.AuthMiddleware()).Use(m.AdminMiddleware())
	admin.POST("/exercises", a.CreateExercise)
	admin.PUT("/exercises/:id", a.UpdateExercise)
	admin.GET("/exercises", a.GetAllExercises)
	admin.DELETE("/exercises/:id", a.DeleteExercise)
}
