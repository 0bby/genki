package db

import (
	"time"

	"github.com/0bby/genki/internal/models"
)

type UpdateApiKeyLabelOpts struct {
	UserID int32
	ID     int32
	Label  string
}

type SaveUserOpts struct {
	Username string
	Password string
	Role     models.UserRole
}

type SaveApiKeyOpts struct {
	Key    string
	UserID int32
	Label  string
}

type UpdateUserOpts struct {
	ID       int32
	Username string
	Password string
}

type GetItemsOpts struct {
	Limit     int
	Page      int
	Timeframe Timeframe
	UserID    int32

	// Filter by exercise
	ExerciseID int
}

type ActivityOpts struct {
	Metric   ActivityMetric
	Step     StepInterval
	Range    int
	Month    int
	Year     int
	Timezone *time.Location
	UserID   int32
}
