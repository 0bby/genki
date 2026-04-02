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
