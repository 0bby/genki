package psql

import (
	"context"
	"fmt"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/models"
)

func (p *Psql) GetExercise(ctx context.Context, id int32) (*models.Exercise, error) {
	row := p.conn.QueryRow(ctx, `SELECT id, name, description, category_id, wger_id, created_at FROM exercises WHERE id = $1`, id)
	e := &models.Exercise{}
	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.CategoryID, &e.WgerID, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("GetExercise: %w", err)
	}
	return e, nil
}

func (p *Psql) GetExerciseByWgerID(ctx context.Context, wgerID int32) (*models.Exercise, error) {
	row := p.conn.QueryRow(ctx, `SELECT id, name, description, category_id, wger_id, created_at FROM exercises WHERE wger_id = $1`, wgerID)
	e := &models.Exercise{}
	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.CategoryID, &e.WgerID, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("GetExerciseByWgerID: %w", err)
	}
	return e, nil
}

func (p *Psql) GetExerciseMuscles(ctx context.Context, exerciseID int32) ([]models.Muscle, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT m.id, m.name, m.name_en, m.is_front, m.wger_id
		FROM muscles m
		JOIN exercise_muscles em ON em.muscle_id = m.id
		WHERE em.exercise_id = $1`, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("GetExerciseMuscles: %w", err)
	}
	defer rows.Close()

	var muscles []models.Muscle
	for rows.Next() {
		var m models.Muscle
		if err := rows.Scan(&m.ID, &m.Name, &m.NameEn, &m.IsFront, &m.WgerID); err != nil {
			return nil, fmt.Errorf("GetExerciseMuscles scan: %w", err)
		}
		muscles = append(muscles, m)
	}
	return muscles, nil
}

func (p *Psql) GetTopExercisesPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[db.RankedItem[*models.Exercise]], error) {
	from, to := opts.Timeframe.Resolve()
	limit := opts.Limit
	if limit <= 0 {
		limit = DefaultItemsPerPage
	}
	offset := (opts.Page - 1) * limit

	rows, err := p.conn.Query(ctx, `
		SELECT e.id, e.name, e.description, e.category_id, e.wger_id, e.created_at,
		       COUNT(ws.id) AS total_sets, COALESCE(SUM(ws.reps), 0) AS total_reps
		FROM exercises e
		JOIN workout_sets ws ON ws.exercise_id = e.id
		JOIN workouts w ON w.id = ws.workout_id
		WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)
		GROUP BY e.id
		ORDER BY total_sets DESC
		LIMIT $4 OFFSET $5`,
		opts.UserID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetTopExercisesPaginated: %w", err)
	}
	defer rows.Close()

	var items []db.RankedItem[*models.Exercise]
	rank := int64(offset + 1)
	for rows.Next() {
		e := &models.Exercise{}
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CategoryID, &e.WgerID, &e.CreatedAt, &e.TotalSets, &e.TotalReps); err != nil {
			return nil, fmt.Errorf("GetTopExercisesPaginated scan: %w", err)
		}
		items = append(items, db.RankedItem[*models.Exercise]{Item: e, Rank: rank})
		rank++
	}

	var totalCount int64
	p.conn.QueryRow(ctx, `
		SELECT COUNT(DISTINCT ws.exercise_id)
		FROM workout_sets ws
		JOIN workouts w ON w.id = ws.workout_id
		WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		opts.UserID, from, to).Scan(&totalCount)

	return &db.PaginatedResponse[db.RankedItem[*models.Exercise]]{
		Items:        items,
		TotalCount:   totalCount,
		ItemsPerPage: int32(limit),
		HasNextPage:  int64(offset+limit) < totalCount,
		CurrentPage:  int32(opts.Page),
	}, nil
}

func (p *Psql) GetTopMuscleGroupsPaginated(ctx context.Context, opts db.GetItemsOpts) (*db.PaginatedResponse[db.RankedItem[*models.Muscle]], error) {
	from, to := opts.Timeframe.Resolve()
	limit := opts.Limit
	if limit <= 0 {
		limit = DefaultItemsPerPage
	}
	offset := (opts.Page - 1) * limit

	rows, err := p.conn.Query(ctx, `
		SELECT m.id, m.name, m.is_front,
		       COUNT(DISTINCT w.id) AS workout_count, COUNT(ws.id) AS set_count
		FROM muscles m
		JOIN exercise_muscles em ON em.muscle_id = m.id
		JOIN workout_sets ws ON ws.exercise_id = em.exercise_id
		JOIN workouts w ON w.id = ws.workout_id
		WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)
		GROUP BY m.id
		ORDER BY set_count DESC
		LIMIT $4 OFFSET $5`,
		opts.UserID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetTopMuscleGroupsPaginated: %w", err)
	}
	defer rows.Close()

	var items []db.RankedItem[*models.Muscle]
	rank := int64(offset + 1)
	for rows.Next() {
		m := &models.Muscle{}
		var workoutCount, setCount int64
		if err := rows.Scan(&m.ID, &m.Name, &m.IsFront, &workoutCount, &setCount); err != nil {
			return nil, fmt.Errorf("GetTopMuscleGroupsPaginated scan: %w", err)
		}
		items = append(items, db.RankedItem[*models.Muscle]{Item: m, Rank: rank})
		rank++
	}

	return &db.PaginatedResponse[db.RankedItem[*models.Muscle]]{
		Items:        items,
		ItemsPerPage: int32(limit),
		HasNextPage:  len(items) == limit,
		CurrentPage:  int32(opts.Page),
	}, nil
}

func (p *Psql) SearchExercises(ctx context.Context, q string) ([]*models.Exercise, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, name, description, category_id, wger_id, created_at
		FROM exercises
		WHERE name % $1 OR name ILIKE '%' || $1 || '%'
		ORDER BY similarity(name, $1) DESC
		LIMIT 20`, q)
	if err != nil {
		return nil, fmt.Errorf("SearchExercises: %w", err)
	}
	defer rows.Close()

	var results []*models.Exercise
	for rows.Next() {
		e := &models.Exercise{}
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.CategoryID, &e.WgerID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("SearchExercises scan: %w", err)
		}
		results = append(results, e)
	}
	return results, nil
}

func (p *Psql) UpsertExercise(ctx context.Context, exercise *models.Exercise) (*models.Exercise, error) {
	row := p.conn.QueryRow(ctx, `
		INSERT INTO exercises (name, description, category_id, wger_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, category_id = EXCLUDED.category_id
		RETURNING id, name, description, category_id, wger_id, created_at`,
		exercise.Name, exercise.Description, exercise.CategoryID, exercise.WgerID)

	e := &models.Exercise{}
	err := row.Scan(&e.ID, &e.Name, &e.Description, &e.CategoryID, &e.WgerID, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("UpsertExercise: %w", err)
	}
	return e, nil
}

func (p *Psql) UpsertExerciseCategory(ctx context.Context, cat *models.ExerciseCategory) (*models.ExerciseCategory, error) {
	row := p.conn.QueryRow(ctx, `
		INSERT INTO exercise_categories (name, wger_id)
		VALUES ($1, $2)
		ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name, wger_id`,
		cat.Name, cat.WgerID)

	c := &models.ExerciseCategory{}
	err := row.Scan(&c.ID, &c.Name, &c.WgerID)
	if err != nil {
		return nil, fmt.Errorf("UpsertExerciseCategory: %w", err)
	}
	return c, nil
}

func (p *Psql) UpsertMuscle(ctx context.Context, m *models.Muscle) (*models.Muscle, error) {
	row := p.conn.QueryRow(ctx, `
		INSERT INTO muscles (name, name_en, is_front, wger_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name, name_en = EXCLUDED.name_en, is_front = EXCLUDED.is_front
		RETURNING id, name, name_en, is_front, wger_id`,
		m.Name, m.NameEn, m.IsFront, m.WgerID)

	result := &models.Muscle{}
	err := row.Scan(&result.ID, &result.Name, &result.NameEn, &result.IsFront, &result.WgerID)
	if err != nil {
		return nil, fmt.Errorf("UpsertMuscle: %w", err)
	}
	return result, nil
}

func (p *Psql) UpsertExerciseMuscle(ctx context.Context, exerciseID, muscleID int32, isPrimary bool) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO exercise_muscles (exercise_id, muscle_id, is_primary)
		VALUES ($1, $2, $3)
		ON CONFLICT (exercise_id, muscle_id) DO UPDATE SET is_primary = EXCLUDED.is_primary`,
		exerciseID, muscleID, isPrimary)
	if err != nil {
		return fmt.Errorf("UpsertExerciseMuscle: %w", err)
	}
	return nil
}
