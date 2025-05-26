package app

import (
	"log"
	middleware "workout-tracker/internal/handler"
	"workout-tracker/internal/handler/admin"
	handler "workout-tracker/internal/handler/auth"
	"workout-tracker/internal/repository/exercise"
	"workout-tracker/internal/repository/user"
	adminService "workout-tracker/internal/service/admin"
	service "workout-tracker/internal/service/auth"
	"workout-tracker/pkg/db"
	"workout-tracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func StartServer() {
	container := dig.New()

	logger.Init("dev")
	err := container.Provide(logger.L)
	if err != nil {
		log.Println("start logger error: ", err)
		return
	}

	err = container.Provide(db.New)
	if err != nil {
		log.Println("start db error: ", err)
		return
	}
	err = container.Provide(user.NewRepository)
	if err != nil {
		log.Println("start user repo error: ", err)
		return
	}
	err = container.Provide(exercise.NewRepository)
	if err != nil {
		log.Println("start exercise repo error: ", err)
		return
	}
	err = container.Provide(service.NewAuthService)
	if err != nil {
		log.Println("start auth service error: ", err)
		return
	}
	err = container.Provide(handler.NewAuthHandler)
	if err != nil {
		log.Println("start auth handler error: ", err)
		return
	}
	err = container.Provide(admin.NewAdminHandler)
	if err != nil {
		log.Println("start admin handler error: ", err)
		return
	}
	err = container.Provide(middleware.NewMiddleware)
	if err != nil {
		log.Println("start middleware error:", err)
		return
	}
	err = container.Provide(adminService.NewAdminService)
	if err != nil {
		log.Println("start admin service error:", err)
		return
	}
	err = container.Provide(gin.Default)
	if err != nil {
		log.Println("start gin error: ", err)
		return
	}

	err = container.Invoke(func(
		router *gin.Engine,
		authHandler *handler.AuthHandler,
		adminHandler *admin.AdminHandler,
		middleware *middleware.Middleware) {
		SetupRoutes(router, authHandler, adminHandler, middleware)
		err := router.Run(":8080")
		if err != nil {
			return
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
