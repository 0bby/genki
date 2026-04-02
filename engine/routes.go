package engine

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/0bby/genki/engine/handlers"
	"github.com/0bby/genki/engine/middleware"
	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	gosync "github.com/0bby/genki/internal/sync"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func bindRoutes(
	r *chi.Mux,
	ready *atomic.Bool,
	db db.DB,
	syncMgr *gosync.Manager,
) {
	if !(len(cfg.AllowedOrigins()) == 0) && !(cfg.AllowedOrigins()[0] == "") {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: cfg.AllowedOrigins(),
			AllowedMethods: []string{"GET", "OPTIONS", "HEAD"},
		}))
	}

	r.Route("/apis/web/v1", func(r chi.Router) {
		r.Get("/config", handlers.GetCfgHandler())

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(db, middleware.AuthModeLoginGate))

			// Activity heatmap
			r.Get("/activity", handlers.GetActivityHandler(db))

			// Workouts
			r.Get("/workouts", handlers.GetWorkoutsHandler(db))
			r.Get("/workout", handlers.GetWorkoutHandler(db))

			// Top lists
			r.Get("/top-exercises", handlers.GetTopExercisesHandler(db))
			r.Get("/top-muscles", handlers.GetTopMusclesHandler(db))

			// Exercise detail
			r.Get("/exercise", handlers.GetExerciseHandler(db))

			// Health data
			r.Get("/sleep", handlers.GetSleepHandler(db))
			r.Get("/heart-rate", handlers.GetHeartRateHandler(db))
			r.Get("/measurements", handlers.GetMeasurementsHandler(db))
			r.Get("/steps", handlers.GetStepsHandler(db))

			// Stats & summary
			r.Get("/stats", handlers.StatsHandler(db))
			r.Get("/summary", handlers.SummaryHandler(db))

			// Search
			r.Get("/search", handlers.SearchHandler(db))

			// Sync status
			r.Get("/sync/status", handlers.SyncStatusHandler(db))
		})

		r.Post("/logout", handlers.LogoutHandler(db))
		if !cfg.RateLimitDisabled() {
			r.With(httprate.Limit(
				10,
				time.Minute,
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
				}),
			)).Post("/login", handlers.LoginHandler(db))
		} else {
			r.Post("/login", handlers.LoginHandler(db))
		}

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			if !ready.Load() {
				http.Error(w, "not ready", http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(db, middleware.AuthModeSessionOrAPIKey))

			// API key management
			r.Get("/user/apikeys", handlers.GetApiKeysHandler(db))
			r.Post("/user/apikeys", handlers.GenerateApiKeyHandler(db))
			r.Patch("/user/apikeys", handlers.UpdateApiKeyLabelHandler(db))
			r.Delete("/user/apikeys", handlers.DeleteApiKeyHandler(db))
			r.Get("/user/me", handlers.MeHandler(db))
			r.Patch("/user", handlers.UpdateUserHandler(db))

			// Manual sync trigger
			r.Post("/sync/trigger", handlers.SyncTriggerHandler(db, syncMgr))

			// Fitbit OAuth
			r.Post("/oauth/fitbit/init", handlers.FitbitOAuthInitHandler())
			r.Get("/oauth/fitbit/callback", handlers.FitbitOAuthCallbackHandler(db))

			// Workout management
			r.Delete("/workout", handlers.DeleteWorkoutHandler(db))
		})
	})

	// serve react client
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "client/build/client"))
	fileServer(r, "/", filesDir)

	// serve client public files
	filesDir = http.Dir(filepath.Join(workDir, "client/public"))
	publicServer(r, "/public", filesDir)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	fs := http.FileServer(root)
	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("client/build/client", strings.TrimPrefix(r.URL.Path, path))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join("client/build/client", "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	})
}

func publicServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}
	fs := http.FileServer(root)
	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, path)
		fs.ServeHTTP(w, r)
	})
}
