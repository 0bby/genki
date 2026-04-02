-- name: UpsertDailySteps :exec
INSERT INTO daily_steps (user_id, date, step_count, source)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, date, source) DO UPDATE SET step_count = EXCLUDED.step_count;

-- name: UpsertDailyActivity :exec
INSERT INTO daily_activity (user_id, date, active_minutes, fairly_active_minutes, lightly_active_minutes, sedentary_minutes, calories_burned, source)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, date, source) DO UPDATE
SET active_minutes = EXCLUDED.active_minutes,
    fairly_active_minutes = EXCLUDED.fairly_active_minutes,
    lightly_active_minutes = EXCLUDED.lightly_active_minutes,
    sedentary_minutes = EXCLUDED.sedentary_minutes,
    calories_burned = EXCLUDED.calories_burned;

-- name: UpsertSleepLog :exec
INSERT INTO sleep_logs (user_id, date, total_minutes, deep_minutes, light_minutes, rem_minutes, awake_minutes, efficiency, start_time, end_time, source, source_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (user_id, date, source) DO UPDATE
SET total_minutes = EXCLUDED.total_minutes,
    deep_minutes = EXCLUDED.deep_minutes,
    light_minutes = EXCLUDED.light_minutes,
    rem_minutes = EXCLUDED.rem_minutes,
    awake_minutes = EXCLUDED.awake_minutes,
    efficiency = EXCLUDED.efficiency,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time;

-- name: UpsertHeartRateDaily :exec
INSERT INTO heart_rate_daily (user_id, date, resting_hr, avg_hr, max_hr, source)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, date, source) DO UPDATE
SET resting_hr = EXCLUDED.resting_hr, avg_hr = EXCLUDED.avg_hr, max_hr = EXCLUDED.max_hr;

-- name: UpsertBodyMeasurement :exec
INSERT INTO body_measurements (user_id, date, weight_kg, body_fat_pct, measurement_category, measurement_value, source)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id, date, measurement_category, source) DO UPDATE
SET weight_kg = EXCLUDED.weight_kg,
    body_fat_pct = EXCLUDED.body_fat_pct,
    measurement_value = EXCLUDED.measurement_value;

-- name: GetDailySteps :many
SELECT * FROM daily_steps
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: GetSleepLogs :many
SELECT * FROM sleep_logs
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: GetHeartRateDaily :many
SELECT * FROM heart_rate_daily
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;

-- name: GetBodyMeasurements :many
SELECT * FROM body_measurements
WHERE user_id = $1 AND date >= $2 AND date <= $3
ORDER BY date DESC;
