package psql

import (
	"context"

	"github.com/0bby/genki/internal/db"
)

func (p *Psql) GetFitnessStats(ctx context.Context, userID int32, tf db.Timeframe) (*db.FitnessStats, error) {
	from, to := tf.Resolve()

	stats := &db.FitnessStats{}

	p.conn.QueryRow(ctx, `SELECT COUNT(*) FROM workouts WHERE user_id = $1 AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)`,
		userID, from, to).Scan(&stats.WorkoutCount)

	p.conn.QueryRow(ctx, `SELECT COUNT(DISTINCT ws.exercise_id) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.ExerciseCount)

	p.conn.QueryRow(ctx, `SELECT COALESCE(SUM(active_minutes), 0) FROM daily_activity WHERE user_id = $1 AND ($2::date IS NULL OR date >= $2::date) AND ($3::date IS NULL OR date < $3::date)`,
		userID, from, to).Scan(&stats.TotalActiveMinutes)

	p.conn.QueryRow(ctx, `SELECT COALESCE(SUM(step_count), 0) FROM daily_steps WHERE user_id = $1 AND ($2::date IS NULL OR date >= $2::date) AND ($3::date IS NULL OR date < $3::date)`,
		userID, from, to).Scan(&stats.TotalSteps)

	p.conn.QueryRow(ctx, `SELECT COALESCE(AVG(total_minutes), 0) FROM sleep_logs WHERE user_id = $1 AND ($2::date IS NULL OR date >= $2::date) AND ($3::date IS NULL OR date < $3::date)`,
		userID, from, to).Scan(&stats.AvgSleepMinutes)

	p.conn.QueryRow(ctx, `SELECT COUNT(*) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.TotalSets)

	p.conn.QueryRow(ctx, `SELECT COALESCE(SUM(ws.reps), 0) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.TotalReps)

	p.conn.QueryRow(ctx, `SELECT COALESCE(AVG(duration_minutes), 0) FROM workouts WHERE user_id = $1 AND duration_minutes IS NOT NULL AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)`,
		userID, from, to).Scan(&stats.AvgWorkoutDuration)

	return stats, nil
}

func (p *Psql) GetRecapStats(ctx context.Context, userID int32, tf db.Timeframe) (*db.RecapStats, error) {
	from, to := tf.Resolve()

	stats := &db.RecapStats{}

	// Basic counts
	p.conn.QueryRow(ctx, `SELECT COUNT(*) FROM workouts WHERE user_id = $1 AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)`,
		userID, from, to).Scan(&stats.TotalWorkouts)

	p.conn.QueryRow(ctx, `SELECT COUNT(*) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.TotalSets)

	p.conn.QueryRow(ctx, `SELECT COALESCE(SUM(ws.reps), 0) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.TotalReps)

	p.conn.QueryRow(ctx, `SELECT COALESCE(SUM(active_minutes), 0) FROM daily_activity WHERE user_id = $1 AND ($2::date IS NULL OR date >= $2::date) AND ($3::date IS NULL OR date < $3::date)`,
		userID, from, to).Scan(&stats.TotalActiveMinutes)

	p.conn.QueryRow(ctx, `SELECT COALESCE(AVG(duration_minutes), 0) FROM workouts WHERE user_id = $1 AND duration_minutes IS NOT NULL AND ($2::timestamptz IS NULL OR started_at >= $2) AND ($3::timestamptz IS NULL OR started_at < $3)`,
		userID, from, to).Scan(&stats.AvgWorkoutDuration)

	p.conn.QueryRow(ctx, `SELECT COUNT(DISTINCT ws.exercise_id) FROM workout_sets ws JOIN workouts w ON w.id = ws.workout_id WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)`,
		userID, from, to).Scan(&stats.ExercisesTried)

	// Top exercises
	rows, err := p.conn.Query(ctx, `
		SELECT e.id, e.name, COUNT(ws.id) AS total_sets, COALESCE(SUM(ws.reps), 0) AS total_reps
		FROM exercises e
		JOIN workout_sets ws ON ws.exercise_id = e.id
		JOIN workouts w ON w.id = ws.workout_id
		WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)
		GROUP BY e.id ORDER BY total_sets DESC LIMIT 10`,
		userID, from, to)
	if err == nil {
		defer rows.Close()
		stats.TopExercises = make([]db.RankedItem[*db.RecapExercise], 0)
		rank := int64(1)
		for rows.Next() {
			e := &db.RecapExercise{}
			if err := rows.Scan(&e.ID, &e.Name, &e.TotalSets, &e.TotalReps); err != nil {
				break
			}
			stats.TopExercises = append(stats.TopExercises, db.RankedItem[*db.RecapExercise]{Item: e, Rank: rank})
			rank++
		}
	}

	// Top muscles
	mRows, err := p.conn.Query(ctx, `
		SELECT m.id, m.name, COALESCE(m.name_en, ''), m.is_front, COUNT(ws.id) AS set_count
		FROM muscles m
		JOIN exercise_muscles em ON em.muscle_id = m.id
		JOIN workout_sets ws ON ws.exercise_id = em.exercise_id
		JOIN workouts w ON w.id = ws.workout_id
		WHERE w.user_id = $1 AND ($2::timestamptz IS NULL OR w.started_at >= $2) AND ($3::timestamptz IS NULL OR w.started_at < $3)
		GROUP BY m.id ORDER BY set_count DESC LIMIT 10`,
		userID, from, to)
	if err == nil {
		defer mRows.Close()
		stats.TopMuscles = make([]db.RankedItem[*db.RecapMuscle], 0)
		rank := int64(1)
		for mRows.Next() {
			m := &db.RecapMuscle{}
			var setCount int64
			if err := mRows.Scan(&m.ID, &m.Name, &m.NameEn, &m.IsFront, &setCount); err != nil {
				break
			}
			stats.TopMuscles = append(stats.TopMuscles, db.RankedItem[*db.RecapMuscle]{Item: m, Rank: rank})
			rank++
		}
	}

	return stats, nil
}
