package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func SearchHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		q := r.URL.Query().Get("q")
		if q == "" {
			http.Error(w, `{"error":"missing query"}`, http.StatusBadRequest)
			return
		}

		exercises, err := store.SearchExercises(r.Context(), q)
		if err != nil {
			l.Error().Err(err).Msg("SearchHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exercises": exercises,
		})
	}
}
