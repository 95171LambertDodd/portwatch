// Package watchdog detects and reports stale or hung scan cycles.
package watchdog

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the liveness state of the watchdog.
type Status string

const (
	StatusOK   Status = "ok"
	StatusStale Status = "stale"
)

// Watchdog tracks the last time a scan was completed and reports staleness.
type Watchdog struct {
	mu       sync.Mutex
	lastBeat time.Time
	maxAge   time.Duration
	clock    func() time.Time
}

// New creates a Watchdog with the given maximum allowed age between heartbeats.
func New(maxAge time.Duration) *Watchdog {
	return &Watchdog{
		maxAge: maxAge,
		clock:  time.Now,
	}
}

func newWithClock(maxAge time.Duration, clock func() time.Time) *Watchdog {
	return &Watchdog{maxAge: maxAge, clock: clock}
}

// Beat records a successful scan tick.
func (w *Watchdog) Beat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.lastBeat = w.clock()
}

// Check returns the current status and a human-readable message.
func (w *Watchdog) Check() (Status, string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.lastBeat.IsZero() {
		return StatusStale, "no heartbeat recorded yet"
	}

	age := w.clock().Sub(w.lastBeat)
	if age > w.maxAge {
		return StatusStale, fmt.Sprintf("last heartbeat was %s ago (max %s)", age.Round(time.Millisecond), w.maxAge)
	}
	return StatusOK, fmt.Sprintf("last heartbeat %s ago", age.Round(time.Millisecond))
}

// LastBeat returns the time of the most recent heartbeat.
func (w *Watchdog) LastBeat() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lastBeat
}
