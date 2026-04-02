package db

import "time"

type ActivityItem struct {
	Start time.Time `json:"start_time"`
	Value int64     `json:"value"`
}

type PaginatedResponse[T any] struct {
	Items        []T   `json:"items"`
	TotalCount   int64 `json:"total_record_count"`
	ItemsPerPage int32 `json:"items_per_page"`
	HasNextPage  bool  `json:"has_next_page"`
	CurrentPage  int32 `json:"current_page"`
}

type RankedItem[T any] struct {
	Item T     `json:"item"`
	Rank int64 `json:"rank"`
}

type FitnessStats struct {
	WorkoutCount       int64   `json:"workout_count"`
	ExerciseCount      int64   `json:"exercise_count"`
	TotalActiveMinutes int64   `json:"total_active_minutes"`
	TotalSteps         int64   `json:"total_steps"`
	AvgSleepMinutes    int64   `json:"avg_sleep_minutes"`
	TotalSets          int64   `json:"total_sets"`
	TotalReps          int64   `json:"total_reps"`
	TotalVolumeKg      float64 `json:"total_volume_kg"`
	AvgWorkoutDuration int64   `json:"avg_workout_duration_min"`
	LongestStreak      int64   `json:"longest_streak"`
	NewExercises       int64   `json:"new_exercises"`
}

type ActivityMetric string

const (
	MetricWorkouts      ActivityMetric = "workouts"
	MetricSteps         ActivityMetric = "steps"
	MetricSleep         ActivityMetric = "sleep"
	MetricActiveMinutes ActivityMetric = "active_minutes"
)
