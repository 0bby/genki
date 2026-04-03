package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"time"

	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/models"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

// WgerSync handles syncing data from a wger instance.
type WgerSync struct {
	store  db.DB
	log    *zerolog.Logger
	client *http.Client
}

func NewWgerSync(store db.DB, log *zerolog.Logger) *WgerSync {
	return &WgerSync{
		store:  store,
		log:    log,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Run starts the periodic sync loop.
func (w *WgerSync) Run(ctx context.Context, stop <-chan struct{}) {
	// Initial sync on startup
	if err := w.SyncAll(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: initial sync failed")
	}

	ticker := time.NewTicker(cfg.WgerSyncInterval())
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			w.log.Info().Msg("wger: sync worker stopped")
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.SyncAll(ctx); err != nil {
				w.log.Error().Err(err).Msg("wger: periodic sync failed")
			}
		}
	}
}

// SyncAll syncs exercises, categories, muscles, workouts, and measurements.
func (w *WgerSync) SyncAll(ctx context.Context) error {
	w.log.Info().Msg("wger: starting full sync")

	if err := w.syncCategories(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: syncCategories failed")
	}
	if err := w.syncMuscles(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: syncMuscles failed")
	}
	if err := w.syncExercises(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: syncExercises failed")
	}
	if err := w.syncWorkoutSessions(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: syncWorkoutSessions failed")
	}
	if err := w.syncWeightEntries(ctx); err != nil {
		w.log.Error().Err(err).Msg("wger: syncWeightEntries failed")
	}

	w.store.UpsertSyncCursor(ctx, &models.SyncCursor{
		UserID:       1,
		Source:       "wger",
		Resource:     "all",
		LastSyncedAt: time.Now(),
	})

	w.log.Info().Msg("wger: sync complete")
	return nil
}

// wger API response types

type wgerPaginatedResponse[T any] struct {
	Count   int    `json:"count"`
	Next    string `json:"next"`
	Results []T    `json:"results"`
}

type wgerCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type wgerMuscle struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	NameEn  string `json:"name_en"`
	IsFront bool   `json:"is_front"`
}

type wgerExercise struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Muscles          []wgerMuscle `json:"muscles"`
	MusclesSecondary []wgerMuscle `json:"muscles_secondary"`
}

type wgerWorkoutSession struct {
	ID         int    `json:"id"`
	Date       string `json:"date"`
	Impression string `json:"impression"` // "1"=general, "2"=good, etc.
	Comment    string `json:"comment"`
	TimeStart  string `json:"time_start"`
	TimeEnd    string `json:"time_end"`
}

type wgerWorkoutLog struct {
	ID            int     `json:"id"`
	Date          string  `json:"date"`
	ExerciseBase  int     `json:"exercise_base"`
	Reps          int     `json:"reps"`
	Weight        string  `json:"weight"`
	WorkoutID     int     `json:"workout"`
	SessionID     *int    `json:"session"`
	RIR           *string `json:"rir"`
}

type wgerWeightEntry struct {
	ID     int    `json:"id"`
	Date   string `json:"date"`
	Weight string `json:"weight"`
}

// API helpers

func (w *WgerSync) get(ctx context.Context, path string) ([]byte, error) {
	url := cfg.WgerURL() + path
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token "+cfg.WgerToken())
	req.Header.Set("Accept", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wger GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wger GET %s: status %d: %s", path, resp.StatusCode, string(body))
	}
	return io.ReadAll(resp.Body)
}

func fetchAll[T any](w *WgerSync, ctx context.Context, path string) ([]T, error) {
	var all []T
	url := path + "?format=json&limit=100&offset=0"
	for url != "" {
		data, err := w.get(ctx, url)
		if err != nil {
			return nil, err
		}
		var page wgerPaginatedResponse[T]
		if err := json.Unmarshal(data, &page); err != nil {
			return nil, fmt.Errorf("wger unmarshal %s: %w", path, err)
		}
		all = append(all, page.Results...)
		if page.Next == "" {
			break
		}
		// Next URL is absolute from wger; extract the path+query
		parsed, err := neturl.Parse(page.Next)
		if err != nil {
			break
		}
		url = parsed.Path
		if parsed.RawQuery != "" {
			url += "?" + parsed.RawQuery
		}
	}
	return all, nil
}

// Sync methods

func (w *WgerSync) syncCategories(ctx context.Context) error {
	cats, err := fetchAll[wgerCategory](w, ctx, "/api/v2/exercisecategory/")
	if err != nil {
		return err
	}
	for _, c := range cats {
		wgerID := int32(c.ID)
		_, err := w.store.UpsertExerciseCategory(ctx, &models.ExerciseCategory{
			Name:   c.Name,
			WgerID: &wgerID,
		})
		if err != nil {
			w.log.Warn().Err(err).Str("name", c.Name).Msg("wger: failed to upsert category")
		}
	}
	w.log.Debug().Int("count", len(cats)).Msg("wger: synced categories")
	return nil
}

func (w *WgerSync) syncMuscles(ctx context.Context) error {
	muscles, err := fetchAll[wgerMuscle](w, ctx, "/api/v2/muscle/")
	if err != nil {
		return err
	}
	for _, m := range muscles {
		wgerID := int32(m.ID)
		_, err := w.store.UpsertMuscle(ctx, &models.Muscle{
			Name:   m.Name,
			NameEn: m.NameEn,
			IsFront: m.IsFront,
			WgerID: &wgerID,
		})
		if err != nil {
			w.log.Warn().Err(err).Str("name", m.Name).Msg("wger: failed to upsert muscle")
		}
	}
	w.log.Debug().Int("count", len(muscles)).Msg("wger: synced muscles")
	return nil
}

func (w *WgerSync) syncExercises(ctx context.Context) error {
	exercises, err := fetchAll[wgerExercise](w, ctx, "/api/v2/exerciseinfo/")
	if err != nil {
		return err
	}
	for _, e := range exercises {
		wgerID := int32(e.ID)

		// Find the local category ID by wger ID
		var catID *int32
		if e.Category.ID > 0 {
			wgerCatID := int32(e.Category.ID)
			// Upsert category inline in case it wasn't synced yet
			cat, err := w.store.UpsertExerciseCategory(ctx, &models.ExerciseCategory{
				Name:   e.Category.Name,
				WgerID: &wgerCatID,
			})
			if err == nil && cat != nil {
				catID = &cat.ID
			}
		}

		ex, err := w.store.UpsertExercise(ctx, &models.Exercise{
			Name:        e.Name,
			Description: e.Description,
			CategoryID:  catID,
			WgerID:      &wgerID,
		})
		if err != nil {
			w.log.Warn().Err(err).Str("name", e.Name).Msg("wger: failed to upsert exercise")
			continue
		}

		// Sync muscle associations
		for _, m := range e.Muscles {
			wgerMID := int32(m.ID)
			muscle, err := w.store.UpsertMuscle(ctx, &models.Muscle{
				Name:   m.Name,
				NameEn: m.NameEn,
				IsFront: m.IsFront,
				WgerID: &wgerMID,
			})
			if err == nil && muscle != nil {
				w.store.UpsertExerciseMuscle(ctx, ex.ID, muscle.ID, true)
			}
		}
		for _, m := range e.MusclesSecondary {
			wgerMID := int32(m.ID)
			muscle, err := w.store.UpsertMuscle(ctx, &models.Muscle{
				Name:   m.Name,
				NameEn: m.NameEn,
				IsFront: m.IsFront,
				WgerID: &wgerMID,
			})
			if err == nil && muscle != nil {
				w.store.UpsertExerciseMuscle(ctx, ex.ID, muscle.ID, false)
			}
		}
	}
	w.log.Debug().Int("count", len(exercises)).Msg("wger: synced exercises")
	return nil
}

func (w *WgerSync) syncWorkoutSessions(ctx context.Context) error {
	// wger uses userID 1 for the single-user self-hosted setup.
	// Genki maps to its own user ID 1 (default admin).
	const userID int32 = 1

	sessions, err := fetchAll[wgerWorkoutSession](w, ctx, "/api/v2/workoutsession/")
	if err != nil {
		return err
	}

	for _, s := range sessions {
		date, err := time.Parse("2006-01-02", s.Date)
		if err != nil {
			w.log.Warn().Str("date", s.Date).Msg("wger: invalid session date")
			continue
		}

		startedAt := date
		if s.TimeStart != "" {
			t, err := time.Parse("15:04:05", s.TimeStart)
			if err == nil {
				startedAt = time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
			}
		}

		var endedAt *time.Time
		if s.TimeEnd != "" {
			t, err := time.Parse("15:04:05", s.TimeEnd)
			if err == nil {
				end := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
				endedAt = &end
			}
		}

		var durationMin *int32
		if endedAt != nil {
			d := int32(endedAt.Sub(startedAt).Minutes())
			durationMin = &d
		}

		sourceID := strconv.Itoa(s.ID)
		workout, err := w.store.SaveWorkout(ctx, &models.Workout{
			UserID:          userID,
			StartedAt:       startedAt,
			EndedAt:         endedAt,
			DurationMinutes: durationMin,
			Title:           s.Comment,
			Source:          "wger",
			SourceID:        &sourceID,
		})
		if err != nil {
			w.log.Warn().Err(err).Int("session_id", s.ID).Msg("wger: failed to save workout")
			continue
		}

		// Sync the logs for this session
		if err := w.syncWorkoutLogs(ctx, workout, s.ID); err != nil {
			w.log.Warn().Err(err).Int("session_id", s.ID).Msg("wger: failed to sync workout logs")
		}
	}

	w.log.Debug().Int("count", len(sessions)).Msg("wger: synced workout sessions")
	return nil
}

func (w *WgerSync) syncWorkoutLogs(ctx context.Context, workout *models.Workout, sessionID int) error {
	path := fmt.Sprintf("/api/v2/workoutlog/?format=json&limit=1000&session=%d", sessionID)
	data, err := w.get(ctx, path)
	if err != nil {
		return err
	}

	var page wgerPaginatedResponse[wgerWorkoutLog]
	if err := json.Unmarshal(data, &page); err != nil {
		return err
	}

	// Delete existing sets and re-insert (idempotent)
	w.store.DeleteWorkoutSetsByWorkout(ctx, workout.ID)

	setNum := int32(0)
	for _, log := range page.Results {
		setNum++

		// Find the exercise by wger ID
		wgerExID := int32(log.ExerciseBase)
		exercise, err := w.store.GetExerciseByWgerID(ctx, wgerExID)
		if err != nil {
			w.log.Debug().Int32("wger_exercise_id", wgerExID).Msg("wger: exercise not found, skipping log")
			continue
		}

		var weightKg *decimal.Decimal
		if log.Weight != "" && log.Weight != "0" {
			d, err := decimal.NewFromString(log.Weight)
			if err == nil {
				weightKg = &d
			}
		}

		var reps *int32
		if log.Reps > 0 {
			r := int32(log.Reps)
			reps = &r
		}

		var rpe *decimal.Decimal
		if log.RIR != nil && *log.RIR != "" {
			// RIR (Reps In Reserve) can be converted to RPE: RPE = 10 - RIR
			rir, err := decimal.NewFromString(*log.RIR)
			if err == nil {
				rpeVal := decimal.NewFromInt(10).Sub(rir)
				rpe = &rpeVal
			}
		}

		_, err = w.store.SaveWorkoutSet(ctx, &models.WorkoutSet{
			WorkoutID:  workout.ID,
			ExerciseID: exercise.ID,
			SetNumber:  setNum,
			Reps:       reps,
			WeightKg:   weightKg,
			RPE:        rpe,
		})
		if err != nil {
			w.log.Warn().Err(err).Int("log_id", log.ID).Msg("wger: failed to save workout set")
		}
	}

	return nil
}

func (w *WgerSync) syncWeightEntries(ctx context.Context) error {
	const userID int32 = 1

	entries, err := fetchAll[wgerWeightEntry](w, ctx, "/api/v2/weightentry/")
	if err != nil {
		return err
	}

	for _, e := range entries {
		date, err := time.Parse("2006-01-02", e.Date)
		if err != nil {
			continue
		}

		weightKg, err := decimal.NewFromString(e.Weight)
		if err != nil {
			continue
		}

		w.store.UpsertBodyMeasurement(ctx, &models.BodyMeasurement{
			UserID:   userID,
			Date:     date,
			WeightKg: &weightKg,
			Source:   "wger",
		})
	}

	w.log.Debug().Int("count", len(entries)).Msg("wger: synced weight entries")
	return nil
}
