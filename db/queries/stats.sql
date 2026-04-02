-- name: CountWorkoutsInPeriod :one
SELECT COUNT(*) FROM workouts
WHERE user_id = $1 AND started_at >= $2 AND started_at < $3;

-- name: CountExercisesInPeriod :one
SELECT COUNT(DISTINCT ws.exercise_id)
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1 AND w.started_at >= $2 AND w.started_at < $3;

-- name: TotalActiveMinutesInPeriod :one
SELECT COALESCE(SUM(active_minutes), 0)::bigint
FROM daily_activity
WHERE user_id = $1 AND date >= $2::date AND date < $3::date;

-- name: TotalStepsInPeriod :one
SELECT COALESCE(SUM(step_count), 0)::bigint
FROM daily_steps
WHERE user_id = $1 AND date >= $2::date AND date < $3::date;

-- name: AvgSleepMinutesInPeriod :one
SELECT COALESCE(AVG(total_minutes), 0)::bigint
FROM sleep_logs
WHERE user_id = $1 AND date >= $2::date AND date < $3::date;

-- name: TotalSetsInPeriod :one
SELECT COUNT(*)
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1 AND w.started_at >= $2 AND w.started_at < $3;

-- name: TotalRepsInPeriod :one
SELECT COALESCE(SUM(ws.reps), 0)::bigint
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1 AND w.started_at >= $2 AND w.started_at < $3;

-- name: TotalVolumeInPeriod :one
SELECT COALESCE(SUM(ws.reps * ws.weight_kg), 0)::numeric
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1 AND w.started_at >= $2 AND w.started_at < $3;

-- name: AvgWorkoutDurationInPeriod :one
SELECT COALESCE(AVG(duration_minutes), 0)::bigint
FROM workouts
WHERE user_id = $1 AND started_at >= $2 AND started_at < $3 AND duration_minutes IS NOT NULL;

-- name: LongestWorkoutStreakInPeriod :one
WITH workout_days AS (
    SELECT DISTINCT date_trunc('day', started_at)::date AS d
    FROM workouts
    WHERE user_id = $1 AND started_at >= $2 AND started_at < $3
),
streaks AS (
    SELECT d, d - (ROW_NUMBER() OVER (ORDER BY d))::int * INTERVAL '1 day' AS grp
    FROM workout_days
)
SELECT COALESCE(MAX(cnt), 0)::bigint AS longest_streak
FROM (SELECT grp, COUNT(*) AS cnt FROM streaks GROUP BY grp) sub;

-- name: NewExercisesInPeriod :one
SELECT COUNT(DISTINCT ws.exercise_id)
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1
  AND w.started_at >= $2
  AND w.started_at < $3
  AND ws.exercise_id NOT IN (
    SELECT DISTINCT ws2.exercise_id
    FROM workout_sets ws2
    JOIN workouts w2 ON w2.id = ws2.workout_id
    WHERE w2.user_id = $1 AND w2.started_at < $2
  );
