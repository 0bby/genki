-- name: UpsertExerciseCategory :one
INSERT INTO exercise_categories (name, wger_id)
VALUES ($1, $2)
ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: UpsertMuscle :one
INSERT INTO muscles (name, name_en, is_front, wger_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name, name_en = EXCLUDED.name_en, is_front = EXCLUDED.is_front
RETURNING *;

-- name: UpsertExercise :one
INSERT INTO exercises (name, description, category_id, wger_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (wger_id) DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description, category_id = EXCLUDED.category_id
RETURNING *;

-- name: UpsertExerciseMuscle :exec
INSERT INTO exercise_muscles (exercise_id, muscle_id, is_primary)
VALUES ($1, $2, $3)
ON CONFLICT (exercise_id, muscle_id) DO UPDATE SET is_primary = EXCLUDED.is_primary;

-- name: GetExercise :one
SELECT * FROM exercises WHERE id = $1;

-- name: GetExerciseByWgerID :one
SELECT * FROM exercises WHERE wger_id = $1;

-- name: GetExerciseMuscles :many
SELECT m.* FROM muscles m
JOIN exercise_muscles em ON em.muscle_id = m.id
WHERE em.exercise_id = $1;

-- name: GetExerciseCategory :one
SELECT * FROM exercise_categories WHERE id = $1;

-- name: GetAllExerciseCategories :many
SELECT * FROM exercise_categories ORDER BY name;

-- name: GetAllMuscles :many
SELECT * FROM muscles ORDER BY name;

-- name: SearchExercises :many
SELECT * FROM exercises
WHERE name % $1 OR name ILIKE '%' || $1 || '%'
ORDER BY similarity(name, $1) DESC
LIMIT 20;

-- name: GetTopExercises :many
SELECT e.id, e.name, e.description, e.category_id, e.wger_id, e.created_at,
       COUNT(ws.id) AS total_sets,
       COALESCE(SUM(ws.reps), 0) AS total_reps
FROM exercises e
JOIN workout_sets ws ON ws.exercise_id = e.id
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1
  AND w.started_at >= $2
  AND w.started_at < $3
GROUP BY e.id
ORDER BY total_sets DESC
LIMIT $4 OFFSET $5;

-- name: CountTopExercises :one
SELECT COUNT(DISTINCT ws.exercise_id)
FROM workout_sets ws
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1
  AND w.started_at >= $2
  AND w.started_at < $3;

-- name: GetTopMuscleGroups :many
SELECT m.id, m.name, m.is_front,
       COUNT(DISTINCT w.id) AS workout_count,
       COUNT(ws.id) AS set_count
FROM muscles m
JOIN exercise_muscles em ON em.muscle_id = m.id
JOIN workout_sets ws ON ws.exercise_id = em.exercise_id
JOIN workouts w ON w.id = ws.workout_id
WHERE w.user_id = $1
  AND w.started_at >= $2
  AND w.started_at < $3
GROUP BY m.id
ORDER BY set_count DESC
LIMIT $4 OFFSET $5;
