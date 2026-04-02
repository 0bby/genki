package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/0bby/genki/engine/middleware"
	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/0bby/genki/internal/db/psql"
	"github.com/0bby/genki/internal/logger"
	"github.com/0bby/genki/internal/models"
	gosync "github.com/0bby/genki/internal/sync"
	"github.com/0bby/genki/internal/utils"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Run(
	getenv func(string) string,
	w io.Writer,
	version string,
) error {
	err := cfg.Load(getenv, version)
	if err != nil {
		panic("Engine: Failed to load configuration")
	}

	l := logger.Get()

	if cfg.StructuredLogging() {
		*l = l.Output(w)
	} else {
		*l = l.Output(zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
			FormatMessage: func(i interface{}) string {
				return fmt.Sprintf("\u001b[30;1m>\u001b[0m %s |", i)
			},
		})
	}

	ctx := logger.NewContext(l)

	l.Info().Msgf("Genki %s", version)

	l.Debug().Msgf("Engine: Checking config directory: %s", cfg.ConfigDir())
	_, err = os.Stat(cfg.ConfigDir())
	if err != nil {
		l.Info().Msgf("Engine: Creating config directory: %s", cfg.ConfigDir())
		err = os.MkdirAll(cfg.ConfigDir(), 0744)
		if err != nil {
			l.Fatal().Err(err).Msg("Engine: Failed to create config directory")
			return err
		}
	}

	l.Debug().Msg("Engine: Initializing database connection")
	var store *psql.Psql
	store, err = psql.New()
	for err != nil {
		l.Error().Err(err).Msg("Engine: Failed to connect to database; retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		store, err = psql.New()
	}
	defer store.Close(ctx)
	l.Info().Msg("Engine: Database connection established")

	if cfg.ForceTZ() != nil {
		l.Debug().Msgf("Engine: Forcing the use of timezone '%s'", cfg.ForceTZ().String())
	}

	l.Debug().Msg("Engine: Checking for default user")
	userCount, _ := store.CountUsers(ctx)
	if userCount < 1 {
		l.Info().Msg("Engine: Creating default user")
		user, err := store.SaveUser(ctx, db.SaveUserOpts{
			Username: cfg.DefaultUsername(),
			Password: cfg.DefaultPassword(),
			Role:     models.UserRoleAdmin,
		})
		if err != nil {
			l.Fatal().Err(err).Msg("Engine: Failed to save default user in database")
		}
		apikey, err := utils.GenerateRandomString(48)
		if err != nil {
			l.Fatal().Err(err).Msg("Engine: Failed to generate default API key")
		}
		_, err = store.SaveApiKey(ctx, db.SaveApiKeyOpts{
			Key:    apikey,
			UserID: user.ID,
			Label:  "Default",
		})
		if err != nil {
			l.Fatal().Err(err).Msg("Engine: Failed to save default API key in database")
		}
		l.Info().Msgf("Engine: Default user created. Login: %s : %s", cfg.DefaultUsername(), cfg.DefaultPassword())
	}

	l.Debug().Msg("Engine: Checking allowed hosts configuration")
	if cfg.AllowAllHosts() {
		l.Warn().Msg("Engine: Configuration allows requests from all hosts. This is a potential security risk!")
	} else if len(cfg.AllowedHosts()) == 0 || cfg.AllowedHosts()[0] == "" {
		l.Warn().Msgf("Engine: No hosts allowed! Did you forget to set the %s variable?", cfg.ALLOWED_HOSTS_ENV)
	} else {
		l.Info().Msgf("Engine: Allowing hosts: %v", cfg.AllowedHosts())
	}

	l.Debug().Msg("Engine: Setting up HTTP server")
	var ready atomic.Bool
	mux := chi.NewRouter()
	mux.Use(middleware.WithRequestID)
	mux.Use(middleware.Logger(l))
	mux.Use(chimiddleware.Recoverer)
	mux.Use(chimiddleware.RealIP)
	mux.Use(middleware.AllowedHosts)
	syncManager := gosync.NewManager(store, l)
	bindRoutes(mux, &ready, store, syncManager)

	httpServer := &http.Server{
		Addr:    cfg.ListenAddr(),
		Handler: mux,
	}

	go func() {
		ready.Store(true)
		l.Info().Msgf("Engine: Listening on %s", cfg.ListenAddr())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("Engine: Error when running ListenAndServe")
		}
	}()

	syncManager.Start(ctx)
	l.Info().Msg("Engine: Initialization finished")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info().Msg("Engine: Received server shutdown notice")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	syncManager.Stop()
	l.Info().Msg("Engine: Waiting for all processes to finish")
	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatal().Err(err).Msg("Engine: Error during server shutdown")
		return err
	}
	l.Info().Msg("Engine: Shutdown successful")
	return nil
}
