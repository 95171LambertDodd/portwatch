package watcher

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alerting"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
)

// Watcher periodically scans ports and emits alerts on changes.
type Watcher struct {
	scanner  *portscanner.Scanner
	alerter  *alerting.Alerter
	cfg      *config.Config
}

// New creates a new Watcher using the provided config.
func New(cfg *config.Config, scanner *portscanner.Scanner, alerter *alerting.Alerter) *Watcher {
	return &Watcher{
		scanner: scanner,
		alerter: alerter,
		cfg:     cfg,
	}
}

// Run starts the watch loop, scanning at the configured interval until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	log.Printf("portwatch: starting watcher (interval=%s)", w.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("portwatch: watcher stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

// tick performs a single scan cycle.
func (w *Watcher) tick() error {
	entries, err := w.scanner.Scan()
	if err != nil {
		return err
	}
	w.alerter.Diff(entries)
	return nil
}
