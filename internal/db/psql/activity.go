package psql

import (
	"context"
	"fmt"

	"github.com/0bby/genki/internal/db"
)

func (p *Psql) GetActivity(ctx context.Context, opts db.ActivityOpts) ([]db.ActivityItem, error) {
	start, end := db.ActivityOptsToTimes(opts)
	interval := stepToInterval(opts.Step)

	var query string
	switch opts.Metric {
	case db.MetricWorkouts:
		query = `
			WITH buckets AS (
				SELECT generate_series($1::timestamptz, $2::timestamptz, $3::interval) AS bucket_start
			)
			SELECT b.bucket_start, COUNT(w.id) AS value
			FROM buckets b
			LEFT JOIN workouts w ON w.started_at >= b.bucket_start AND w.started_at < b.bucket_start + $3::interval AND w.user_id = $4
			GROUP BY b.bucket_start ORDER BY b.bucket_start`
	case db.MetricSteps:
		query = `
			WITH buckets AS (
				SELECT generate_series($1::timestamptz, $2::timestamptz, $3::interval) AS bucket_start
			)
			SELECT b.bucket_start, COALESCE(SUM(ds.step_count), 0)::bigint AS value
			FROM buckets b
			LEFT JOIN daily_steps ds ON ds.date >= b.bucket_start::date AND ds.date < (b.bucket_start + $3::interval)::date AND ds.user_id = $4
			GROUP BY b.bucket_start ORDER BY b.bucket_start`
	case db.MetricSleep:
		query = `
			WITH buckets AS (
				SELECT generate_series($1::timestamptz, $2::timestamptz, $3::interval) AS bucket_start
			)
			SELECT b.bucket_start, COALESCE(SUM(sl.total_minutes), 0)::bigint AS value
			FROM buckets b
			LEFT JOIN sleep_logs sl ON sl.date >= b.bucket_start::date AND sl.date < (b.bucket_start + $3::interval)::date AND sl.user_id = $4
			GROUP BY b.bucket_start ORDER BY b.bucket_start`
	case db.MetricActiveMinutes:
		query = `
			WITH buckets AS (
				SELECT generate_series($1::timestamptz, $2::timestamptz, $3::interval) AS bucket_start
			)
			SELECT b.bucket_start, COALESCE(SUM(da.active_minutes), 0)::bigint AS value
			FROM buckets b
			LEFT JOIN daily_activity da ON da.date >= b.bucket_start::date AND da.date < (b.bucket_start + $3::interval)::date AND da.user_id = $4
			GROUP BY b.bucket_start ORDER BY b.bucket_start`
	default:
		return nil, fmt.Errorf("GetActivity: unknown metric %s", opts.Metric)
	}

	rows, err := p.conn.Query(ctx, query, start, end, interval, opts.UserID)
	if err != nil {
		return nil, fmt.Errorf("GetActivity: %w", err)
	}
	defer rows.Close()

	var items []db.ActivityItem
	for rows.Next() {
		var item db.ActivityItem
		if err := rows.Scan(&item.Start, &item.Value); err != nil {
			return nil, fmt.Errorf("GetActivity scan: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}
