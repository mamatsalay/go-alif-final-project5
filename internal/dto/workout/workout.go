package workout

import (
	"workout-tracker/internal/model/workout"
	"workout-tracker/internal/model/workoutexercisejoin"
)

type CreateWorkoutWithExercisesRequest struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	Exercises []struct {
		ExerciseID int `json:"exercise_id"`
		Sets       int `json:"sets"`
		Reps       int `json:"reps"`
	} `json:"exercises"`
}

type CreateWorkoutWithExercisesResponse struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	Category  string `json:"category"`
	ID        int    `json:"id"`
	Exercises []struct {
		ExerciseName string `json:"exercise_name"`
		Sets         int    `json:"sets"`
		Reps         int    `json:"reps"`
	} `json:"exercises"`
}

type CreateWorkoutRequest struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Category string `json:"category"`
	UserID   int    `json:"user_id"`
}

type UpdateWorkoutRequest struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Category string `json:"category"`
	UserID   int    `json:"user_id"`
	ID       int    `json:"id"`
}

type WorkoutWithExercises struct {
	workout.Workout
	Exercises []workoutexercisejoin.WorkoutExercise `json:"exercises"`
}
