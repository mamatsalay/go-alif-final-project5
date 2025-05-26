package workoutexercisejoin

import "workout-tracker/internal/model/exercise"

type WorkoutExercise struct {
	WorkoutID  int                `json:"workout_id"`
	ExerciseID int                `json:"exercise_id"`
	Reps       int                `json:"reps,omitempty"`
	Sets       int                `json:"sets,omitempty"`
	Exercise   *exercise.Exercise `json:"exercise,omitempty"`
}
