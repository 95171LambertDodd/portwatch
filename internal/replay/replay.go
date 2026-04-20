// Package replay provides a mechanism to replay historical port events
// through the alert pipeline, useful for testing rules against past data
// or re-processing events after configuration changes.
package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/notify"
)

// Event represents a single replayed history event passed to the handler.
type Event struct {
	Timestamp time.Time
	Port      uint16
	Protocol  string
	Address   string
	PID       int
	Process   string
	Tags      []string
}

// Handler is called for each event during replay.
type Handler func(ctx context.Context, e Event) error

// Replayer reads events from a history recorder and re-emits them.
type Replayer struct {
	recorder *history.Recorder
	notifier *notify.Notifier
	delay    time.Duration
}

// Config holds options for constructing a Replayer.
type Config struct {
	// Recorder is the source of historical events.
	Recorder *history.Recorder
	// Notifier receives replayed events as notifications.
	Notifier *notify.Notifier
	// Delay is an optional pause between replayed events (0 = no delay).
	Delay time.Duration
}

// New constructs a Replayer from the given Config.
// Returns an error if required fields are missing.
func New(cfg Config) (*Replayer, error) {
	if cfg.Recorder == nil {
		return nil, fmt.Errorf("replay: recorder is required")
	}
	if cfg.Notifier == nil {
		return nil, fmt.Errorf("replay: notifier is required")
	}
	return &Replayer{
		recorder: cfg.Recorder,
		notifier: cfg.Notifier,
		delay:    cfg.Delay,
	}, nil
}

// Run reads all recorded events and re-emits them through the notifier.
// It respects context cancellation between events.
func (r *Replayer) Run(ctx context.Context) error {
	events, err := r.recorder.ReadAll()
	if err != nil {
		return fmt.Errorf("replay: reading history: %w", err)
	}

	for _, raw := range events {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := r.notifier.Notify(ctx, raw); err != nil {
			return fmt.Errorf("replay: notifying event: %w", err)
		}

		if r.delay > 0 {
			select {
			case <-time.After(r.delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return nil
}

// Count returns the number of events available for replay without
// actually replaying them.
func (r *Replayer) Count() (int, error) {
	events, err := r.recorder.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("replay: reading history for count: %w", err)
	}
	return len(events), nil
}
