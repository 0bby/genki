package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/logger"
	gosync "github.com/0bby/genki/internal/sync"
)

func FitbitOAuthInitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.FitbitEnabled() {
			http.Error(w, `{"error":"fitbit not configured"}`, http.StatusBadRequest)
			return
		}

		stateBytes := make([]byte, 16)
		rand.Read(stateBytes)
		state := hex.EncodeToString(stateBytes)

		// Store state in a cookie for CSRF protection
		http.SetCookie(w, &http.Cookie{
			Name:     "fitbit_oauth_state",
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})

		url := gosync.FitbitAuthURL(state)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"url": url,
		})
	}
}

func FitbitOAuthCallbackHandler(store db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.FromContext(r.Context())

		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" {
			http.Error(w, `{"error":"missing code"}`, http.StatusBadRequest)
			return
		}

		// Verify CSRF state
		cookie, err := r.Cookie("fitbit_oauth_state")
		if err != nil || cookie.Value != state {
			http.Error(w, `{"error":"invalid state"}`, http.StatusBadRequest)
			return
		}

		// Clear the state cookie
		http.SetCookie(w, &http.Cookie{
			Name:   "fitbit_oauth_state",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		userID := getUserFromContext(r)
		fitbit := gosync.NewFitbitSync(store, logger.Get())

		if err := fitbit.ExchangeCode(r.Context(), userID, code); err != nil {
			l.Error().Err(err).Msg("FitbitOAuthCallback: token exchange failed")
			http.Error(w, `{"error":"token exchange failed"}`, http.StatusInternalServerError)
			return
		}

		// Trigger an initial sync in the background (use background context since request will end)
		go func() {
			if err := fitbit.SyncUser(context.Background(), userID); err != nil {
				l.Error().Err(err).Int32("user_id", userID).Msg("FitbitOAuthCallback: initial sync failed")
			}
		}()

		// Redirect back to home page
		http.Redirect(w, r, "/?fitbit=connected", http.StatusFound)
	}
}
