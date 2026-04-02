package sync

import (
	"context"

	"github.com/0bby/genki/internal/cfg"
	"github.com/0bby/genki/internal/db"
	"github.com/rs/zerolog"
)

// Manager coordinates all sync workers.
type Manager struct {
	store db.DB
	log   *zerolog.Logger
	wger  *WgerSync
	stop  chan struct{}
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
	return m
}

// Start launches all configured sync workers in the background.
func (m *Manager) Start(ctx context.Context) {
	if m.wger != nil {
		m.log.Info().Msg("Sync: Starting wger sync worker")
		go m.wger.Run(ctx, m.stop)
	}
	// Fitbit sync is per-user and triggered after OAuth; not a global worker.
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
