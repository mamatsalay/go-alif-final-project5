package admin

import (
	"net/http"
	"strconv"
	dto "workout-tracker/internal/dto/exercise"
	"workout-tracker/internal/erorrs"
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

func (h *AdminHandler) CreateExercise(c *gin.Context) {
	var ex dto.CreateExerciseRequest
	if err := c.ShouldBindJSON(&ex); err != nil {
		h.Logger.Errorw("Invalid exercise data", erorrs.ErrorKey, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	created, err := h.Service.CreateExercise(c.Request.Context(), ex)
	if err != nil {
		h.Logger.Errorw("Failed to create exercise", erorrs.ErrorKey, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"exercise_id": created})
}

// UpdateExercise обрабатывает запрос PUT /admin/exercises/:id (обновляет занятие).
func (h *AdminHandler) UpdateExercise(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Errorw("Invalid exercise ID", "id", idStr, erorrs.ErrorKey, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var ex dto.CreateExerciseRequest
	if err := c.ShouldBindJSON(&ex); err != nil {
		h.Logger.Errorw("Invalid exercise data", erorrs.ErrorKey, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	err = h.Service.UpdateExercise(c.Request.Context(), id, ex)
	if err != nil {
		h.Logger.Errorw("Failed to update exercise", "id", id, erorrs.ErrorKey, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"answer": "exercise successfully updated"})
}

func (h *AdminHandler) GetAllExercises(c *gin.Context) {
	exercises, err := h.Service.GetAllExercises(c.Request.Context())
	if err != nil {
		h.Logger.Errorw("Failed to retrieve exercises", erorrs.ErrorKey, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercises"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *AdminHandler) DeleteExercise(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Errorw("Invalid exercise ID", "id", idStr, erorrs.ErrorKey, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	err = h.Service.DeleteExercise(c.Request.Context(), id)
	if err != nil {
		h.Logger.Errorw("Failed to delete exercise", erorrs.ErrorKey, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"answer": "exercise successfully deleted"})
}
