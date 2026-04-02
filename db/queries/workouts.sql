-- name: InsertWorkout :one
INSERT INTO workouts (user_id, started_at, ended_at, duration_minutes, title, notes, source, source_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (user_id, source, source_id) DO UPDATE
SET started_at = EXCLUDED.started_at,
    ended_at = EXCLUDED.ended_at,
    duration_minutes = EXCLUDED.duration_minutes,
    title = EXCLUDED.title,
    notes = EXCLUDED.notes
RETURNING *;

-- name: InsertWorkoutSet :one
INSERT INTO workout_sets (workout_id, exercise_id, set_number, reps, weight_kg, duration_seconds, rpe)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetWorkout :one
SELECT * FROM workouts WHERE id = $1;

-- name: GetWorkoutSets :many
SELECT ws.*, e.name AS exercise_name
FROM workout_sets ws
JOIN exercises e ON e.id = ws.exercise_id
WHERE ws.workout_id = $1
ORDER BY ws.set_number;

-- name: GetWorkoutsPaginated :many
SELECT * FROM workouts
WHERE user_id = $1
  AND started_at >= $2
  AND started_at < $3
ORDER BY started_at DESC
LIMIT $4 OFFSET $5;

-- name: CountWorkouts :one
SELECT COUNT(*) FROM workouts
WHERE user_id = $1
  AND started_at >= $2
  AND started_at < $3;

-- name: DeleteWorkoutSetsByWorkout :exec
DELETE FROM workout_sets WHERE workout_id = $1;

-- name: DeleteWorkout :exec
DELETE FROM workouts WHERE id = $1 AND user_id = $2;
