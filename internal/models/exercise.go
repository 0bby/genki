package models

import "time"

type ExerciseCategory struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	WgerID *int32 `json:"wger_id,omitempty"`
}

type Muscle struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	NameEn  string `json:"name_en,omitempty"`
	IsFront bool   `json:"is_front"`
	WgerID  *int32 `json:"wger_id,omitempty"`
}

type Exercise struct {
	ID          int32             `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	CategoryID  *int32            `json:"category_id,omitempty"`
	Category    *ExerciseCategory `json:"category,omitempty"`
	WgerID      *int32            `json:"wger_id,omitempty"`
	Muscles     []Muscle          `json:"muscles,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`

	// Aggregated stats (populated by queries)
	TotalSets    int64   `json:"total_sets,omitempty"`
	TotalReps    int64   `json:"total_reps,omitempty"`
	TotalVolumeKg float64 `json:"total_volume_kg,omitempty"`
}
