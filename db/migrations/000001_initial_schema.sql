-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- Extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

-- Types
CREATE TYPE role AS ENUM (
    'admin',
    'user'
);

-- ============================================================
-- Auth tables (preserved from Koito)
-- ============================================================

CREATE TABLE users (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY (
        SEQUENCE NAME users_id_seq
        START WITH 1
        INCREMENT BY 1
        NO MINVALUE
        NO MAXVALUE
        CACHE 1
    ),
    username text UNIQUE NOT NULL,
    role role DEFAULT 'user'::role NOT NULL,
    password bytea NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE api_keys (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY (
        SEQUENCE NAME api_keys_id_seq
        START WITH 1
        INCREMENT BY 1
        NO MINVALUE
        NO MAXVALUE
        CACHE 1
    ),
    key text UNIQUE NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    label text NOT NULL,
    CONSTRAINT api_keys_pkey PRIMARY KEY (id)
);

CREATE TABLE sessions (
    id UUID NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    persistent boolean DEFAULT false NOT NULL,
    CONSTRAINT sessions_pkey PRIMARY KEY (id)
);

-- ============================================================
-- Exercise domain
-- ============================================================

CREATE TABLE exercise_categories (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE,
    wger_id integer UNIQUE
);

CREATE TABLE muscles (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE,
    name_en text,
    is_front boolean NOT NULL DEFAULT true,
    wger_id integer UNIQUE
);

CREATE TABLE exercises (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL,
    description text DEFAULT '',
    category_id integer REFERENCES exercise_categories(id),
    wger_id integer UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE exercise_muscles (
    exercise_id integer NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    muscle_id integer NOT NULL REFERENCES muscles(id) ON DELETE CASCADE,
    is_primary boolean NOT NULL DEFAULT true,
    PRIMARY KEY (exercise_id, muscle_id)
);

-- ============================================================
-- Workout tracking
-- ============================================================

CREATE TABLE workouts (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    started_at timestamptz NOT NULL,
    ended_at timestamptz,
    duration_minutes integer,
    title text DEFAULT '',
    notes text DEFAULT '',
    source text NOT NULL DEFAULT 'manual',
    source_id text,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE(user_id, source, source_id)
);

CREATE TABLE workout_sets (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    workout_id integer NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id integer NOT NULL REFERENCES exercises(id),
    set_number integer NOT NULL DEFAULT 1,
    reps integer,
    weight_kg numeric(7,2),
    duration_seconds integer,
    rpe numeric(3,1),
    logged_at timestamptz NOT NULL DEFAULT now()
);

-- ============================================================
-- Passive / biometric data (Fitbit, etc.)
-- ============================================================

CREATE TABLE daily_steps (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date date NOT NULL,
    step_count integer NOT NULL,
    source text NOT NULL DEFAULT 'fitbit',
    UNIQUE(user_id, date, source)
);

CREATE TABLE daily_activity (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date date NOT NULL,
    active_minutes integer NOT NULL DEFAULT 0,
    fairly_active_minutes integer DEFAULT 0,
    lightly_active_minutes integer DEFAULT 0,
    sedentary_minutes integer DEFAULT 0,
    calories_burned integer,
    source text NOT NULL DEFAULT 'fitbit',
    UNIQUE(user_id, date, source)
);

CREATE TABLE sleep_logs (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date date NOT NULL,
    total_minutes integer NOT NULL,
    deep_minutes integer,
    light_minutes integer,
    rem_minutes integer,
    awake_minutes integer,
    efficiency integer,
    start_time timestamptz,
    end_time timestamptz,
    source text NOT NULL DEFAULT 'fitbit',
    source_id text,
    UNIQUE(user_id, date, source)
);

CREATE TABLE heart_rate_daily (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date date NOT NULL,
    resting_hr integer,
    avg_hr integer,
    max_hr integer,
    source text NOT NULL DEFAULT 'fitbit',
    UNIQUE(user_id, date, source)
);

CREATE TABLE body_measurements (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date date NOT NULL,
    weight_kg numeric(6,2),
    body_fat_pct numeric(5,2),
    measurement_category text,
    measurement_value numeric(8,2),
    source text NOT NULL DEFAULT 'wger',
    UNIQUE(user_id, date, measurement_category, source)
);

-- ============================================================
-- Sync infrastructure
-- ============================================================

CREATE TABLE sync_cursors (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source text NOT NULL,
    resource text NOT NULL,
    last_synced_at timestamptz NOT NULL DEFAULT now(),
    cursor_value text,
    UNIQUE(user_id, source, resource)
);

CREATE TABLE oauth_tokens (
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider text NOT NULL,
    access_token text NOT NULL,
    refresh_token text NOT NULL,
    expires_at timestamptz NOT NULL,
    scopes text,
    UNIQUE(user_id, provider)
);

-- ============================================================
-- Foreign key constraints (auth tables)
-- ============================================================

ALTER TABLE ONLY api_keys
    ADD CONSTRAINT api_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ONLY sessions
    ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- ============================================================
-- Indexes
-- ============================================================

-- Exercise search
CREATE INDEX idx_exercises_name_trgm ON exercises USING gin (name gin_trgm_ops);
CREATE INDEX idx_exercises_category ON exercises(category_id);

-- Workout queries
CREATE INDEX idx_workouts_user_started ON workouts(user_id, started_at);
CREATE INDEX idx_workouts_source ON workouts(source);
CREATE INDEX idx_workout_sets_workout ON workout_sets(workout_id);
CREATE INDEX idx_workout_sets_exercise ON workout_sets(exercise_id);

-- Daily data queries
CREATE INDEX idx_daily_steps_user_date ON daily_steps(user_id, date);
CREATE INDEX idx_daily_activity_user_date ON daily_activity(user_id, date);
CREATE INDEX idx_sleep_logs_user_date ON sleep_logs(user_id, date);
CREATE INDEX idx_heart_rate_user_date ON heart_rate_daily(user_id, date);
CREATE INDEX idx_body_measurements_user_date ON body_measurements(user_id, date);

-- Sync
CREATE INDEX idx_sync_cursors_user_source ON sync_cursors(user_id, source);
CREATE INDEX idx_oauth_tokens_user ON oauth_tokens(user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS oauth_tokens CASCADE;
DROP TABLE IF EXISTS sync_cursors CASCADE;
DROP TABLE IF EXISTS body_measurements CASCADE;
DROP TABLE IF EXISTS heart_rate_daily CASCADE;
DROP TABLE IF EXISTS sleep_logs CASCADE;
DROP TABLE IF EXISTS daily_activity CASCADE;
DROP TABLE IF EXISTS daily_steps CASCADE;
DROP TABLE IF EXISTS workout_sets CASCADE;
DROP TABLE IF EXISTS workouts CASCADE;
DROP TABLE IF EXISTS exercise_muscles CASCADE;
DROP TABLE IF EXISTS exercises CASCADE;
DROP TABLE IF EXISTS muscles CASCADE;
DROP TABLE IF EXISTS exercise_categories CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop types
DROP TYPE IF EXISTS role;

-- Drop extensions
DROP EXTENSION IF EXISTS pg_trgm;
