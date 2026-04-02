package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailySteps struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Date      time.Time `json:"date"`
	StepCount int32     `json:"step_count"`
	Source    string    `json:"source"`
}

type DailyActivity struct {
	ID                    int32     `json:"id"`
	UserID                int32     `json:"user_id"`
	Date                  time.Time `json:"date"`
	ActiveMinutes         int32     `json:"active_minutes"`
	FairlyActiveMinutes   int32     `json:"fairly_active_minutes"`
	LightlyActiveMinutes  int32     `json:"lightly_active_minutes"`
	SedentaryMinutes      int32     `json:"sedentary_minutes"`
	CaloriesBurned        *int32    `json:"calories_burned,omitempty"`
	Source                string    `json:"source"`
}

type SleepLog struct {
	ID           int32      `json:"id"`
	UserID       int32      `json:"user_id"`
	Date         time.Time  `json:"date"`
	TotalMinutes int32      `json:"total_minutes"`
	DeepMinutes  *int32     `json:"deep_minutes,omitempty"`
	LightMinutes *int32     `json:"light_minutes,omitempty"`
	REMMinutes   *int32     `json:"rem_minutes,omitempty"`
	AwakeMinutes *int32     `json:"awake_minutes,omitempty"`
	Efficiency   *int32     `json:"efficiency,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Source       string     `json:"source"`
	SourceID     *string    `json:"source_id,omitempty"`
}

type HeartRateDaily struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Date      time.Time `json:"date"`
	RestingHR *int32    `json:"resting_hr,omitempty"`
	AvgHR     *int32    `json:"avg_hr,omitempty"`
	MaxHR     *int32    `json:"max_hr,omitempty"`
	Source    string    `json:"source"`
}

type BodyMeasurement struct {
	ID                  int32            `json:"id"`
	UserID              int32            `json:"user_id"`
	Date                time.Time        `json:"date"`
	WeightKg            *decimal.Decimal `json:"weight_kg,omitempty"`
	BodyFatPct          *decimal.Decimal `json:"body_fat_pct,omitempty"`
	MeasurementCategory *string          `json:"measurement_category,omitempty"`
	MeasurementValue    *decimal.Decimal `json:"measurement_value,omitempty"`
	Source              string           `json:"source"`
}
