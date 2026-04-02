-- name: WorkoutActivity :many
WITH buckets AS (
    SELECT generate_series($1::timestamptz, $2::timestamptz, $3::interval) AS bucket_start
),
bucketed AS (
    SELECT b.bucket_start, COUNT(w.id) AS value
    FROM buckets b
    LEFT JOIN workouts w
        ON w.started_at >= b.bucket_start
        AND w.started_at < b.bucket_start + $3::interval
        AND w.user_id = $4
    GROUP BY b.bucket_start
    ORDER BY b.bucket_start
)
SELECT * FROM bucketed;

-- name: StepsActivity :many
WITH buckets AS (
    SELECT generate_series($1::date, $2::date, $3::interval) AS bucket_start
),
bucketed AS (
    SELECT b.bucket_start, COALESCE(SUM(ds.step_count), 0)::bigint AS value
    FROM buckets b
    LEFT JOIN daily_steps ds
        ON ds.date >= b.bucket_start::date
        AND ds.date < (b.bucket_start + $3::interval)::date
        AND ds.user_id = $4
    GROUP BY b.bucket_start
    ORDER BY b.bucket_start
)
SELECT * FROM bucketed;

-- name: SleepActivity :many
WITH buckets AS (
    SELECT generate_series($1::date, $2::date, $3::interval) AS bucket_start
),
bucketed AS (
    SELECT b.bucket_start, COALESCE(SUM(sl.total_minutes), 0)::bigint AS value
    FROM buckets b
    LEFT JOIN sleep_logs sl
        ON sl.date >= b.bucket_start::date
        AND sl.date < (b.bucket_start + $3::interval)::date
        AND sl.user_id = $4
    GROUP BY b.bucket_start
    ORDER BY b.bucket_start
)
SELECT * FROM bucketed;

-- name: ActiveMinutesActivity :many
WITH buckets AS (
    SELECT generate_series($1::date, $2::date, $3::interval) AS bucket_start
),
bucketed AS (
    SELECT b.bucket_start, COALESCE(SUM(da.active_minutes), 0)::bigint AS value
    FROM buckets b
    LEFT JOIN daily_activity da
        ON da.date >= b.bucket_start::date
        AND da.date < (b.bucket_start + $3::interval)::date
        AND da.user_id = $4
    GROUP BY b.bucket_start
    ORDER BY b.bucket_start
)
SELECT * FROM bucketed;
