package statistics

type WorkoutStatistics struct {
	Category    string  `json:"category"`
	TotalWeight float64 `json:"total_weight"`
	UserID      int     `json:"user_id"`
	TotalSets   int     `json:"total_sets"`
	TotalReps   int     `json:"total_reps"`
}
