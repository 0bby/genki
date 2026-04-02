package psql

import (
	"context"
	"fmt"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/models"
)

func (p *Psql) GetWorkout(ctx context.Context, id int32) (*models.Workout, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT id, user_id, started_at, ended_at, duration_minutes, title, notes, source, source_id, created_at
		FROM workouts WHERE id = $1`, id)

	w := &models.Workout{}
	err := row.Scan(&w.ID, &w.UserID, &w.StartedAt, &w.EndedAt, &w.DurationMinutes, &w.Title, &w.Notes, &w.Source, &w.SourceID, &w.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("GetWorkout: %w", err)
	}
	return w, nil
}

func (p *Psql) GetWorkoutSets(ctx context.Context, workoutID int32) ([]models.WorkoutSet, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT ws.id, ws.workout_id, ws.exercise_id, ws.set_number, ws.reps, ws.weight_kg, ws.duration_seconds, ws.rpe, ws.logged_at
		FROM workout_sets ws
		WHERE ws.workout_id = $1
		ORDER BY ws.set_number`, workoutID)
	if err != nil {
		return nil, fmt.Errorf("GetWorkoutSets: %w", err)
	}
	defer rows.Close()

	var sets []models.WorkoutSet
	for rows.Next() {
		var s models.WorkoutSet
		if err := rows.Scan(&s.ID, &s.WorkoutID, &s.ExerciseID, &s.SetNumber, &s.Reps, &s.WeightKg, &s.DurationSeconds, &s.RPE, &s.LoggedAt); err != nil {
			return nil, fmt.Errorf("GetWorkoutSets scan: %w", err)
		}
		sets = append(sets, s)
	}
	return sets, nil
}

func (p *Psql) GetWorkoutsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[*models.Workout], error) {
	from, to := opts.Timeframe.Resolve()
	limit := opts.Limit
	if limit <= 0 {
		limit = DefaultItemsPerPage
	}
	offset := (opts.Page - 1) * limit

	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, started_at, ended_at, duration_minutes, title, notes, source, source_id, created_at
		FROM workouts
		WHERE user_id = $1 AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)
		ORDER BY started_at DESC
		LIMIT $4 OFFSET $5`,
		opts.UserID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetWorkoutsPaginated: %w", err)
	}
	defer rows.Close()

	var items []*models.Workout
	for rows.Next() {
		w := &models.Workout{}
		if err := rows.Scan(&w.ID, &w.UserID, &w.StartedAt, &w.EndedAt, &w.DurationMinutes, &w.Title, &w.Notes, &w.Source, &w.SourceID, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetWorkoutsPaginated scan: %w", err)
		}
		items = append(items, w)
	}

	var totalCount int64
	p.conn.QueryRow(ctx, `
		SELECT COUNT(*) FROM workouts
		WHERE user_id = $1 AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)`,
		opts.UserID, from, to).Scan(&totalCount)

	return &db.PaginatedResponse[*models.Workout]{
		Items:        items,
		TotalCount:   totalCount,
		ItemsPerPage: int32(limit),
		HasNextPage:  int64(offset+limit) < totalCount,
		CurrentPage:  int32(opts.Page),
	}, nil
}

func (p *Psql) SaveWorkout(ctx context.Context, w *models.Workout) (*models.Workout, error) {
	row := p.conn.QueryRow(ctx, `
		INSERT INTO workouts (user_id, started_at, ended_at, duration_minutes, title, notes, source, source_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, source, source_id) DO UPDATE
		SET started_at = EXCLUDED.started_at, ended_at = EXCLUDED.ended_at,
		    duration_minutes = EXCLUDED.duration_minutes, title = EXCLUDED.title, notes = EXCLUDED.notes
		RETURNING id, user_id, started_at, ended_at, duration_minutes, title, notes, source, source_id, created_at`,
		w.UserID, w.StartedAt, w.EndedAt, w.DurationMinutes, w.Title, w.Notes, w.Source, w.SourceID)

	result := &models.Workout{}
	err := row.Scan(&result.ID, &result.UserID, &result.StartedAt, &result.EndedAt, &result.DurationMinutes, &result.Title, &result.Notes, &result.Source, &result.SourceID, &result.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("SaveWorkout: %w", err)
	}
	return result, nil
}

func (p *Psql) SaveWorkoutSet(ctx context.Context, s *models.WorkoutSet) (*models.WorkoutSet, error) {
	row := p.conn.QueryRow(ctx, `
		INSERT INTO workout_sets (workout_id, exercise_id, set_number, reps, weight_kg, duration_seconds, rpe)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, workout_id, exercise_id, set_number, reps, weight_kg, duration_seconds, rpe, logged_at`,
		s.WorkoutID, s.ExerciseID, s.SetNumber, s.Reps, s.WeightKg, s.DurationSeconds, s.RPE)

	result := &models.WorkoutSet{}
	err := row.Scan(&result.ID, &result.WorkoutID, &result.ExerciseID, &result.SetNumber, &result.Reps, &result.WeightKg, &result.DurationSeconds, &result.RPE, &result.LoggedAt)
	if err != nil {
		return nil, fmt.Errorf("SaveWorkoutSet: %w", err)
	}
	return result, nil
}

func (p *Psql) DeleteWorkoutSetsByWorkout(ctx context.Context, workoutID int32) error {
	_, err := p.conn.Exec(ctx, `DELETE FROM workout_sets WHERE workout_id = $1`, workoutID)
	return err
}

func (p *Psql) DeleteWorkout(ctx context.Context, id, userID int32) error {
	_, err := p.conn.Exec(ctx, `DELETE FROM workouts WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
