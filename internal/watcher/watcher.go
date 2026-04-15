package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher orchestrates periodic port scanning, diffing, alerting, and snapshotting.
type Watcher struct {
	cfg     config.Config
	scanner *portscanner.Scanner
	alerter *alerting.Alerter
	store   *snapshot.Store
}

// New constructs a Watcher from the provided configuration.
func New(cfg config.Config) *Watcher {
	return &Watcher{
		cfg:     cfg,
		scanner: portscanner.NewScanner(),
		alerter: alerting.NewAlerter(),
		store:   snapshot.NewStore(cfg.SnapshotPath),
	}
}

// Run starts the watch loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	// Perform an immediate tick on startup.
	if err := w.tick(); err != nil {
		log.Printf("[portwatch] initial scan error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := w.tick(); err != nil {
				log.Printf("[portwatch] scan error: %v", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *Watcher) tick() error {
	prev, err := w.store.Load()
	if err != nil {
		return err
	}

	current, err := w.scanner.Scan()
	if err != nil {
		return err
	}

	alerts := w.alerter.Diff(prev.Entries, current)
	for _, a := range alerts {
		w.alerter.Emit(a)
	}

	return w.store.Save(snapshot.Snapshot{
		CapturedAt: time.Now().UTC(),
		Entries:    current,
	})
}
