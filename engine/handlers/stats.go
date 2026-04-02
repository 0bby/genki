package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func StatsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		tf := TimeframeFromRequest(r)

		stats, err := store.GetFitnessStats(r.Context(), userID, tf)
		if err != nil {
			l.Error().Err(err).Msg("StatsHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}

func SummaryHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		tf := TimeframeFromRequest(r)

		stats, err := store.GetFitnessStats(r.Context(), userID, tf)
		if err != nil {
			l.Error().Err(err).Msg("SummaryHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// TODO: Expand with top exercises, muscle groups, etc.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
