package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func GetTopExercisesHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		opts := OptsFromRequest(r, userID)

		result, err := store.GetTopExercisesPaginated(r.Context(), opts)
		if err != nil {
			l.Error().Err(err).Msg("GetTopExercisesHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func GetTopMusclesHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		opts := OptsFromRequest(r, userID)

		result, err := store.GetTopMuscleGroupsPaginated(r.Context(), opts)
		if err != nil {
			l.Error().Err(err).Msg("GetTopMusclesHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
