package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func timeRangeFromRequest(r *http.Request) (time.Time, time.Time) {
	tf := TimeframeFromRequest(r)
	from, to := tf.Resolve()
	return from, to
}

func GetSleepHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		from, to := timeRangeFromRequest(r)

		logs, err := store.GetSleepLogs(r.Context(), userID, from, to)
		if err != nil {
			l.Error().Err(err).Msg("GetSleepHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

func GetHeartRateHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		from, to := timeRangeFromRequest(r)

		data, err := store.GetHeartRateDaily(r.Context(), userID, from, to)
		if err != nil {
			l.Error().Err(err).Msg("GetHeartRateHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func GetMeasurementsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		from, to := timeRangeFromRequest(r)

		data, err := store.GetBodyMeasurements(r.Context(), userID, from, to)
		if err != nil {
			l.Error().Err(err).Msg("GetMeasurementsHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}

func GetStepsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		from, to := timeRangeFromRequest(r)

		data, err := store.GetDailySteps(r.Context(), userID, from, to)
		if err != nil {
			l.Error().Err(err).Msg("GetStepsHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	}
}
