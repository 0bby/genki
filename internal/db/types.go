package db

import "time"

type ActivityItem struct {
	Start time.Time `json:"start"`
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
	AvgWorkoutDuration int64   `json:"avg_workout_duration"`
	LongestStreak      int64   `json:"longest_streak"`
	NewExercises       int64   `json:"new_exercises"`
}

type RecapStats struct {
	Title              string                      `json:"title"`
	TopExercises       []RankedItem[*RecapExercise] `json:"top_exercises"`
	TopMuscles         []RankedItem[*RecapMuscle]   `json:"top_muscles"`
	TotalWorkouts      int64                        `json:"total_workouts"`
	TotalSets          int64                        `json:"total_sets"`
	TotalReps          int64                        `json:"total_reps"`
	TotalActiveMinutes int64                        `json:"total_active_minutes"`
	AvgWorkoutDuration float64                      `json:"avg_workout_duration"`
	ExercisesTried     int64                        `json:"exercises_tried"`
	NewExercises       int64                        `json:"new_exercises"`
	WorkoutStreak      int64                        `json:"workout_streak"`
}

type RecapExercise struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	TotalSets  int64  `json:"total_sets"`
	TotalReps  int64  `json:"total_reps"`
}

type RecapMuscle struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	NameEn  string `json:"name_en"`
	IsFront bool   `json:"is_front"`
}

type ActivityMetric string

const (
	MetricWorkouts      ActivityMetric = "workouts"
	MetricSteps         ActivityMetric = "steps"
	MetricSleep         ActivityMetric = "sleep"
	MetricActiveMinutes ActivityMetric = "active_minutes"
)
