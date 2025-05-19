package exercise

import "time"

type Exercise struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"update_at"`
	Name        string    `json:"name"`
	ID          int       `json:"id"`
	WorkoutID   int       `json:"workout_id"`
	Sets        int       `json:"sets"`
	Reps        int       `json:"reps"`
	Weight      float64   `json:"weight"`
	Description string    `json:"description"`
}
