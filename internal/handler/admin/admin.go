// internal/handler/admin/admin.go
package admin

import (
	"net/http"
	"strconv"
	"workout-tracker/internal/model/exercise"
	"workout-tracker/internal/service/admin"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

type AdminHandlerParams struct {
	dig.In
	Service *admin.AdminService
	Logger  *zap.SugaredLogger
}

type AdminHandler struct {
	Service *admin.AdminService
	Logger  *zap.SugaredLogger
}

func NewAdminHandler(p AdminHandlerParams) *AdminHandler {
	return &AdminHandler{
		Service: p.Service,
		Logger:  p.Logger,
	}
}

// CreateExercise обрабатывает запрос POST /admin/exercises (создаёт занятие)
func (h *AdminHandler) CreateExercise(c *gin.Context) {
	var ex exercise.Exercise
	if err := c.ShouldBindJSON(&ex); err != nil {
		h.Logger.Errorw("Invalid exercise data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	created, err := h.Service.CreateExercise(c.Request.Context(), ex)
	if err != nil {
		h.Logger.Errorw("Failed to create exercise", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// UpdateExercise обрабатывает запрос PUT /admin/exercises/:id (обновляет занятие)
func (h *AdminHandler) UpdateExercise(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Errorw("Invalid exercise ID", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var ex exercise.Exercise
	if err := c.ShouldBindJSON(&ex); err != nil {
		h.Logger.Errorw("Invalid exercise data", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}
	ex.ID = id

	updated, err := h.Service.UpdateExercise(c.Request.Context(), ex)
	if err != nil {
		h.Logger.Errorw("Failed to update exercise", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// GetAllExercises обрабатывает запрос GET /admin/exercises (выводит все занятия)
func (h *AdminHandler) GetAllExercises(c *gin.Context) {
	exercises, err := h.Service.GetAllExercises(c.Request.Context())
	if err != nil {
		h.Logger.Errorw("Failed to retrieve exercises", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercises"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}
