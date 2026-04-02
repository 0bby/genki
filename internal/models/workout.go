package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Workout struct {
	ID              int32        `json:"id"`
	UserID          int32        `json:"user_id"`
	StartedAt       time.Time    `json:"started_at"`
	EndedAt         *time.Time   `json:"ended_at,omitempty"`
	DurationMinutes *int32       `json:"duration_minutes,omitempty"`
	Title           string       `json:"title,omitempty"`
	Notes           string       `json:"notes,omitempty"`
	Source          string       `json:"source"`
	SourceID        *string      `json:"source_id,omitempty"`
	CreatedAt       time.Time    `json:"created_at"`
	Sets            []WorkoutSet `json:"sets,omitempty"`
}

type WorkoutSet struct {
	ID              int32            `json:"id"`
	WorkoutID       int32            `json:"workout_id"`
	ExerciseID      int32            `json:"exercise_id"`
	Exercise        *Exercise        `json:"exercise,omitempty"`
	SetNumber       int32            `json:"set_number"`
	Reps            *int32           `json:"reps,omitempty"`
	WeightKg        *decimal.Decimal `json:"weight_kg,omitempty"`
	DurationSeconds *int32           `json:"duration_seconds,omitempty"`
	RPE             *decimal.Decimal `json:"rpe,omitempty"`
	LoggedAt        time.Time        `json:"logged_at"`
}
