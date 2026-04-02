package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func GetExerciseHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		exercise, err := store.GetExercise(r.Context(), int32(id))
		if err != nil {
			l.Error().Err(err).Msg("GetExerciseHandler: failed")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		muscles, err := store.GetExerciseMuscles(r.Context(), int32(id))
		if err != nil {
			l.Error().Err(err).Msg("GetExerciseHandler: failed to get muscles")
		}
		exercise.Muscles = muscles

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(exercise)
	}
}
