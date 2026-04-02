package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/0bby/genki/internal/models"
)

func (p *Psql) UpsertDailySteps(ctx context.Context, s *models.DailySteps) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO daily_steps (user_id, date, step_count, source)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, date, source) DO UPDATE SET step_count = EXCLUDED.step_count`,
		s.UserID, s.Date, s.StepCount, s.Source)
	return err
}

func (p *Psql) UpsertDailyActivity(ctx context.Context, a *models.DailyActivity) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO daily_activity (user_id, date, active_minutes, fairly_active_minutes, lightly_active_minutes, sedentary_minutes, calories_burned, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, date, source) DO UPDATE
		SET active_minutes = EXCLUDED.active_minutes, fairly_active_minutes = EXCLUDED.fairly_active_minutes,
		    lightly_active_minutes = EXCLUDED.lightly_active_minutes, sedentary_minutes = EXCLUDED.sedentary_minutes,
		    calories_burned = EXCLUDED.calories_burned`,
		a.UserID, a.Date, a.ActiveMinutes, a.FairlyActiveMinutes, a.LightlyActiveMinutes, a.SedentaryMinutes, a.CaloriesBurned, a.Source)
	return err
}

func (p *Psql) UpsertSleepLog(ctx context.Context, s *models.SleepLog) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO sleep_logs (user_id, date, total_minutes, deep_minutes, light_minutes, rem_minutes, awake_minutes, efficiency, start_time, end_time, source, source_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (user_id, date, source) DO UPDATE
		SET total_minutes = EXCLUDED.total_minutes, deep_minutes = EXCLUDED.deep_minutes,
		    light_minutes = EXCLUDED.light_minutes, rem_minutes = EXCLUDED.rem_minutes,
		    awake_minutes = EXCLUDED.awake_minutes, efficiency = EXCLUDED.efficiency,
		    start_time = EXCLUDED.start_time, end_time = EXCLUDED.end_time`,
		s.UserID, s.Date, s.TotalMinutes, s.DeepMinutes, s.LightMinutes, s.REMMinutes, s.AwakeMinutes, s.Efficiency, s.StartTime, s.EndTime, s.Source, s.SourceID)
	return err
}

func (p *Psql) UpsertHeartRateDaily(ctx context.Context, h *models.HeartRateDaily) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO heart_rate_daily (user_id, date, resting_hr, avg_hr, max_hr, source)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, date, source) DO UPDATE
		SET resting_hr = EXCLUDED.resting_hr, avg_hr = EXCLUDED.avg_hr, max_hr = EXCLUDED.max_hr`,
		h.UserID, h.Date, h.RestingHR, h.AvgHR, h.MaxHR, h.Source)
	return err
}

func (p *Psql) UpsertBodyMeasurement(ctx context.Context, m *models.BodyMeasurement) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO body_measurements (user_id, date, weight_kg, body_fat_pct, measurement_category, measurement_value, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, date, measurement_category, source) DO UPDATE
		SET weight_kg = EXCLUDED.weight_kg, body_fat_pct = EXCLUDED.body_fat_pct, measurement_value = EXCLUDED.measurement_value`,
		m.UserID, m.Date, m.WeightKg, m.BodyFatPct, m.MeasurementCategory, m.MeasurementValue, m.Source)
	return err
}

func (p *Psql) GetDailySteps(ctx context.Context, userID int32, from, to time.Time) ([]models.DailySteps, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, date, step_count, source FROM daily_steps
		WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetDailySteps: %w", err)
	}
	defer rows.Close()
	var result []models.DailySteps
	for rows.Next() {
		var s models.DailySteps
		if err := rows.Scan(&s.ID, &s.UserID, &s.Date, &s.StepCount, &s.Source); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (p *Psql) GetSleepLogs(ctx context.Context, userID int32, from, to time.Time) ([]models.SleepLog, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, date, total_minutes, deep_minutes, light_minutes, rem_minutes, awake_minutes, efficiency, start_time, end_time, source, source_id
		FROM sleep_logs WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetSleepLogs: %w", err)
	}
	defer rows.Close()
	var result []models.SleepLog
	for rows.Next() {
		var s models.SleepLog
		if err := rows.Scan(&s.ID, &s.UserID, &s.Date, &s.TotalMinutes, &s.DeepMinutes, &s.LightMinutes, &s.REMMinutes, &s.AwakeMinutes, &s.Efficiency, &s.StartTime, &s.EndTime, &s.Source, &s.SourceID); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (p *Psql) GetHeartRateDaily(ctx context.Context, userID int32, from, to time.Time) ([]models.HeartRateDaily, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, date, resting_hr, avg_hr, max_hr, source
		FROM heart_rate_daily WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetHeartRateDaily: %w", err)
	}
	defer rows.Close()
	var result []models.HeartRateDaily
	for rows.Next() {
		var h models.HeartRateDaily
		if err := rows.Scan(&h.ID, &h.UserID, &h.Date, &h.RestingHR, &h.AvgHR, &h.MaxHR, &h.Source); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (p *Psql) GetBodyMeasurements(ctx context.Context, userID int32, from, to time.Time) ([]models.BodyMeasurement, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, date, weight_kg, body_fat_pct, measurement_category, measurement_value, source
		FROM body_measurements WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC`, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("GetBodyMeasurements: %w", err)
	}
	defer rows.Close()
	var result []models.BodyMeasurement
	for rows.Next() {
		var m models.BodyMeasurement
		if err := rows.Scan(&m.ID, &m.UserID, &m.Date, &m.WeightKg, &m.BodyFatPct, &m.MeasurementCategory, &m.MeasurementValue, &m.Source); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}
