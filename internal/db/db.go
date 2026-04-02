// package db defines the database interface
package db

import (
	"context"
	"time"

	"github.com/0bby/genki/internal/models"
	"github.com/google/uuid"
)

type DB interface {
	// ============================================================
	// Auth (preserved from Koito)
	// ============================================================

	GetUserBySession(ctx context.Context, sessionId uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByApiKey(ctx context.Context, key string) (*models.User, error)
	GetApiKeysByUserID(ctx context.Context, id int32) ([]models.ApiKey, error)
	SaveUser(ctx context.Context, opts SaveUserOpts) (*models.User, error)
	SaveApiKey(ctx context.Context, opts SaveApiKeyOpts) (*models.ApiKey, error)
	SaveSession(ctx context.Context, userId int32, expiresAt time.Time, persistent bool) (*models.Session, error)
	UpdateUser(ctx context.Context, opts UpdateUserOpts) error
	UpdateApiKeyLabel(ctx context.Context, opts UpdateApiKeyLabelOpts) error
	RefreshSession(ctx context.Context, sessionId uuid.UUID, expiresAt time.Time) error
	DeleteSession(ctx context.Context, sessionId uuid.UUID) error
	DeleteApiKey(ctx context.Context, id int32) error
	CountUsers(ctx context.Context) (int64, error)

	// ============================================================
	// Exercises
	// ============================================================

	GetExercise(ctx context.Context, id int32) (*models.Exercise, error)
	GetExerciseByWgerID(ctx context.Context, wgerID int32) (*models.Exercise, error)
	GetExerciseMuscles(ctx context.Context, exerciseID int32) ([]models.Muscle, error)
	GetTopExercisesPaginated(ctx context.Context, opts GetItemsOpts) (*PaginatedResponse[RankedItem[*models.Exercise]], error)
	GetTopMuscleGroupsPaginated(ctx context.Context, opts GetItemsOpts) (*PaginatedResponse[RankedItem[*models.Muscle]], error)
	SearchExercises(ctx context.Context, q string) ([]*models.Exercise, error)
	UpsertExercise(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error)
	UpsertExerciseCategory(ctx context.Context, cat *models.ExerciseCategory) (*models.ExerciseCategory, error)
	UpsertMuscle(ctx context.Context, m *models.Muscle) (*models.Muscle, error)
	UpsertExerciseMuscle(ctx context.Context, exerciseID, muscleID int32, isPrimary bool) error

	// ============================================================
	// Workouts
	// ============================================================

	GetWorkout(ctx context.Context, id int32) (*models.Workout, error)
	GetWorkoutSets(ctx context.Context, workoutID int32) ([]models.WorkoutSet, error)
	GetWorkoutsPaginated(ctx context.Context, opts GetItemsOpts) (*PaginatedResponse[*models.Workout], error)
	SaveWorkout(ctx context.Context, w *models.Workout) (*models.Workout, error)
	SaveWorkoutSet(ctx context.Context, s *models.WorkoutSet) (*models.WorkoutSet, error)
	DeleteWorkoutSetsByWorkout(ctx context.Context, workoutID int32) error
	DeleteWorkout(ctx context.Context, id, userID int32) error

	// ============================================================
	// Health data
	// ============================================================

	UpsertDailySteps(ctx context.Context, s *models.DailySteps) error
	UpsertDailyActivity(ctx context.Context, a *models.DailyActivity) error
	UpsertSleepLog(ctx context.Context, s *models.SleepLog) error
	UpsertHeartRateDaily(ctx context.Context, h *models.HeartRateDaily) error
	UpsertBodyMeasurement(ctx context.Context, m *models.BodyMeasurement) error
	GetDailySteps(ctx context.Context, userID int32, from, to time.Time) ([]models.DailySteps, error)
	GetSleepLogs(ctx context.Context, userID int32, from, to time.Time) ([]models.SleepLog, error)
	GetHeartRateDaily(ctx context.Context, userID int32, from, to time.Time) ([]models.HeartRateDaily, error)
	GetBodyMeasurements(ctx context.Context, userID int32, from, to time.Time) ([]models.BodyMeasurement, error)

	// ============================================================
	// Activity heatmap
	// ============================================================

	GetActivity(ctx context.Context, opts ActivityOpts) ([]ActivityItem, error)

	// ============================================================
	// Stats
	// ============================================================

	GetFitnessStats(ctx context.Context, userID int32, tf Timeframe) (*FitnessStats, error)

	// ============================================================
	// Sync
	// ============================================================

	UpsertSyncCursor(ctx context.Context, c *models.SyncCursor) error
	GetSyncCursor(ctx context.Context, userID int32, source, resource string) (*models.SyncCursor, error)
	GetAllSyncCursors(ctx context.Context, userID int32) ([]models.SyncCursor, error)
	UpsertOAuthToken(ctx context.Context, t *models.OAuthToken) error
	GetOAuthToken(ctx context.Context, userID int32, provider string) (*models.OAuthToken, error)
	DeleteOAuthToken(ctx context.Context, userID int32, provider string) error

	// ============================================================
	// System
	// ============================================================

	Ping(ctx context.Context) error
	Close(ctx context.Context)
}
