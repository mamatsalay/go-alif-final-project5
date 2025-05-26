package workout

import (
	"net/http"
	"strconv"
	dto "workout-tracker/internal/dto/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
	"workout-tracker/internal/service/workout"

	"github.com/gin-gonic/gin"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type WorkoutHandlerParams struct {
	dig.In

	Service *workout.WorkoutService
	Logger  *zap.SugaredLogger
}

type WorkoutHandler struct {
	Service *workout.WorkoutService
	Log     *zap.SugaredLogger
}

func NewWorkoutHandler(params WorkoutHandlerParams) *WorkoutHandler {
	return &WorkoutHandler{
		Service: params.Service,
		Log:     params.Logger,
	}
}

func (h *WorkoutHandler) Create(c *gin.Context) {
	var req dto.CreateWorkoutWithExercisesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetInt("userID")

	var exercises []workoutexercisejoin.WorkoutExercise
	for _, e := range req.Exercises {
		exercises = append(exercises, workoutexercisejoin.WorkoutExercise{
			ExerciseID: e.ExerciseID,
			Reps:       e.Reps,
			Sets:       e.Sets,
		})
	}

	if err := h.Service.CreateWorkout(c.Request.Context(), userID, req.Name, req.Title, req.Category, exercises); err != nil {
		h.Log.Errorw("error creating workout", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create workout"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "workout created"})
}

func (h *WorkoutHandler) Update(c *gin.Context) {
	var req dto.CreateWorkoutWithExercisesRequest
	workoutID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercises []workoutexercisejoin.WorkoutExercise
	for _, e := range req.Exercises {
		exercises = append(exercises, workoutexercisejoin.WorkoutExercise{
			ExerciseID: e.ExerciseID,
			Reps:       e.Reps,
			Sets:       e.Sets,
		})
	}

	if err := h.Service.UpdateWorkout(c, userID, workoutID, req.Name, req.Title, req.Category, exercises); err != nil {
		h.Log.Errorw("error updating workout", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update workout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workout updated"})
}

func (h *WorkoutHandler) Delete(c *gin.Context) {
	workoutID, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("userID")

	if err := h.Service.DeleteWorkout(c, userID, workoutID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete workout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workout deleted"})
}

func (h *WorkoutHandler) GetAll(c *gin.Context) {
	userID := c.GetInt("userID")

	workouts, err := h.Service.GetAllWorkoutsWithExercises(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get workouts"})
		return
	}

	c.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) Get(c *gin.Context) {
	userID := c.GetInt("userID")
	workoutIDStr := c.Param("id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	result, err := h.Service.GetWorkoutByID(c, userID, workoutID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "workout not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, result)
}
