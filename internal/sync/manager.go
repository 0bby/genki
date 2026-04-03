package sync

import (
	"context"
	"time"

	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/rs/zerolog"
)

// Manager coordinates all sync workers.
type Manager struct {
	store  db.DB
	log    *zerolog.Logger
	wger   *WgerSync
	fitbit *FitbitSync
	stop   chan struct{}
}

func NewManager(store db.DB, log *zerolog.Logger) *Manager {
	m := &Manager{
		store: store,
		log:   log,
		stop:  make(chan struct{}),
	}
	if cfg.WgerEnabled() {
		m.wger = NewWgerSync(store, log)
	}
	if cfg.FitbitEnabled() {
		m.fitbit = NewFitbitSync(store, log)
	}
	return m
}

// Start launches all configured sync workers in the background.
func (m *Manager) Start(ctx context.Context) {
	if m.wger != nil {
		m.log.Info().Msg("Sync: Starting wger sync worker")
		go m.wger.Run(ctx, m.stop)
	}
	if m.fitbit != nil {
		m.log.Info().Msg("Sync: Starting Fitbit sync worker")
		go m.runFitbitWorker(ctx)
	}
}

// runFitbitWorker periodically syncs Fitbit data for user 1.
func (m *Manager) runFitbitWorker(ctx context.Context) {
	const userID int32 = 1

	// Initial sync — only if a token exists (user has connected Fitbit)
	if err := m.fitbit.SyncUser(ctx, userID); err != nil {
		m.log.Debug().Err(err).Msg("Sync: Fitbit initial sync skipped (no token or error)")
	}

	ticker := time.NewTicker(cfg.FitbitSyncInterval())
	defer ticker.Stop()

	for {
		select {
		case <-m.stop:
			m.log.Info().Msg("Sync: Fitbit worker stopped")
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := m.fitbit.SyncUser(ctx, userID); err != nil {
				m.log.Warn().Err(err).Msg("Sync: Fitbit periodic sync failed")
			}
		}
	}
}

// Stop signals all workers to shut down.
func (m *Manager) Stop() {
	close(m.stop)
}

// TriggerWgerSync runs a one-off wger sync if configured.
func (m *Manager) TriggerWgerSync(ctx context.Context) error {
	if m.wger == nil {
		return nil
	}
	return m.wger.SyncAll(ctx)
}

// TriggerFitbitSync runs a one-off Fitbit sync if configured.
func (m *Manager) TriggerFitbitSync(ctx context.Context, userID int32) error {
	if m.fitbit == nil {
		return nil
	}
	return m.fitbit.SyncUser(ctx, userID)
}
