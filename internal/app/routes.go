package app

import (
	"workout-tracker/internal/handler"
	"workout-tracker/internal/handler/admin"
	"workout-tracker/internal/handler/auth"
	"workout-tracker/internal/handler/workout"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *auth.AuthHandler, a *admin.AdminHandler, w *workout.WorkoutHandler, m *handler.Middleware) {
	auth := r.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)

	admin := r.Group("/admin").Use(m.AuthMiddleware()).Use(m.AdminMiddleware())
	admin.POST("/exercises", a.CreateExercise)
	admin.PUT("/exercises/:id", a.UpdateExercise)
	admin.GET("/exercises", a.GetAllExercises)
	admin.DELETE("/exercises/:id", a.DeleteExercise)

	workout := r.Group("/workouts").Use(m.AuthMiddleware())
	workout.POST("", w.Create)
	workout.PUT("/:id", w.Update)
	workout.GET("", w.GetAll)
	workout.GET("/:id", w.Get)
	workout.DELETE("/:id", w.Delete)
	workout.POST("/:id/photo", w.UpdatePhoto)
	workout.GET("/:id/photo", w.GetPhoto)

	allExercises := r.Group("/exercises").Use(m.AuthMiddleware())
	allExercises.GET("", a.GetAllExercises)
}
