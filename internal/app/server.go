package app

import (
	"log"
	"workout-tracker/internal/handler/admin"
	handler "workout-tracker/internal/handler/auth"
	"workout-tracker/internal/repository/exercise"
	"workout-tracker/internal/repository/user"
	service "workout-tracker/internal/service/auth"
	"workout-tracker/pkg/db"
	"workout-tracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

func StartServer() {
	container := dig.New()

	logger.Init("dev")
	err := container.Provide(logger.Sugar)
	if err != nil {
		return
	}

	err = container.Provide(db.New)
	if err != nil {
		return
	}
	err = container.Provide(user.NewRepository)
	if err != nil {
		return
	}
	err = container.Provide(exercise.NewRepository)
	if err != nil {
		return
	}
	err = container.Provide(service.NewAuthService)
	if err != nil {
		return
	}
	err = container.Provide(handler.NewAuthHandler)
	if err != nil {
		return
	}
	err = container.Provide(admin.NewAdminHandler)
	if err != nil {
		return
	}
	err = container.Provide(gin.Default)
	if err != nil {
		return
	}

	err = container.Invoke(func(
		router *gin.Engine,
		authHandler *handler.AuthHandler,
		adminHandler *admin.AdminHandler,
		serviceHandler *service.AuthService) {
		SetupRoutes(router, authHandler, adminHandler, serviceHandler)
		err := router.Run(":8080")
		if err != nil {
			return
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}
