package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
	gosync "github.com/0bby/genki/internal/sync"
)

func SyncStatusHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		userID := getUserFromContext(r)

		cursors, err := store.GetAllSyncCursors(r.Context(), userID)
		if err != nil {
			l.Error().Err(err).Msg("SyncStatusHandler: failed")
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sources": cursors,
		})
	}
}

func SyncTriggerHandler(store db.DB, syncMgr *gosync.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())
		source := r.URL.Query().Get("source")

		w.Header().Set("Content-Type", "application/json")

		switch source {
		case "wger":
			go func() {
				if err := syncMgr.TriggerWgerSync(context.Background()); err != nil {
					l.Error().Err(err).Msg("SyncTriggerHandler: wger sync failed")
				}
			}()
			json.NewEncoder(w).Encode(map[string]string{
				"status": "triggered",
				"source": "wger",
			})
		case "fitbit":
			userID := getUserFromContext(r)
			fitbit := gosync.NewFitbitSync(store, logger.Get())
			go func() {
				if err := fitbit.SyncUser(context.Background(), userID); err != nil {
					l.Error().Err(err).Int32("user_id", userID).Msg("SyncTriggerHandler: fitbit sync failed")
				}
			}()
			json.NewEncoder(w).Encode(map[string]string{
				"status": "triggered",
				"source": "fitbit",
			})
		default:
			http.Error(w, `{"error":"unknown source"}`, http.StatusBadRequest)
		}
	}
}
