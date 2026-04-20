// Package circuitbreaker implements a simple circuit breaker to prevent
// repeated alerting or processing when a downstream sink or scanner is
// consistently failing.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // failing, requests blocked
	StateHalfOpen              // testing if recovery is possible
)

// ErrOpen is returned when the circuit breaker is open.
var ErrOpen = errors.New("circuit breaker is open")

// Breaker is a circuit breaker that trips after a threshold of consecutive
// failures and resets after a cooldown period.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	threshold    int
	cooldown     time.Duration
	lastFailure  time.Time
	clock        func() time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts recovery after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	return newWithClock(threshold, cooldown, time.Now)
}

func newWithClock(threshold int, cooldown time.Duration, clock func() time.Time) *Breaker {
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
		clock:     clock,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit
// is open. It transitions to half-open once the cooldown has elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if b.clock().Sub(b.lastFailure) >= b.cooldown {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the breaker to closed state.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure count and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.lastFailure = b.clock()
	if b.failures >= b.threshold {
		b.state = StateOpen
	}
}

// State returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
