package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
)

func GetWorkoutsHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)
		opts := OptsFromRequest(r, userID)

		result, err := store.GetWorkoutsPaginated(r.Context(), opts)
		if err != nil {
			l.Error().Err(err).Msg("GetWorkoutsHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

func GetWorkoutHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		workout, err := store.GetWorkout(r.Context(), int32(id))
		if err != nil {
			l.Error().Err(err).Msg("GetWorkoutHandler: failed to get workout")
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		sets, err := store.GetWorkoutSets(r.Context(), int32(id))
		if err != nil {
			l.Error().Err(err).Msg("GetWorkoutHandler: failed to get sets")
		}
		workout.Sets = sets

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workout)
	}
}

func DeleteWorkoutHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id < 1 {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		err = store.DeleteWorkout(r.Context(), int32(id), userID)
		if err != nil {
			l.Error().Err(err).Msg("DeleteWorkoutHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
