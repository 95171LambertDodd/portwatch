package circuitbreaker

import (
	"testing"
	"time"
)

var baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

type fixedClock struct {
	now time.Time
}

func (f *fixedClock) Time() time.Time { return f.now }

func newTestBreaker(threshold int, cooldown time.Duration, clk *fixedClock) *Breaker {
	return newWithClock(threshold, cooldown, clk.Time)
}

func TestNew_ReturnsNonNil(t *testing.T) {
	b := New(3, time.Second)
	if b == nil {
		t.Fatal("expected non-nil Breaker")
	}
}

func TestAllow_InitiallyAllowed(t *testing.T) {
	clk := &fixedClock{now: baseTime}
	b := newTestBreaker(3, time.Second, clk)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	clk := &fixedClock{now: baseTime}
	b := newTestBreaker(3, 10*time.Second, clk)

	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}

	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	clk := &fixedClock{now: baseTime}
	b := newTestBreaker(2, 5*time.Second, clk)

	b.RecordFailure()
	b.RecordFailure()

	// still within cooldown
	clk.now = baseTime.Add(3 * time.Second)
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen within cooldown, got %v", err)
	}

	// advance past cooldown
	clk.now = baseTime.Add(6 * time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil after cooldown, got %v", err)
	}
	if b.CurrentState() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.CurrentState())
	}
}

func TestRecordSuccess_ResetsToClosed(t *testing.T) {
	clk := &fixedClock{now: baseTime}
	b := newTestBreaker(2, 5*time.Second, clk)

	b.RecordFailure()
	b.RecordFailure()

	// advance past cooldown to half-open
	clk.now = baseTime.Add(6 * time.Second)
	_ = b.Allow()

	b.RecordSuccess()
	if b.CurrentState() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", b.CurrentState())
	}
}

func TestRecordFailure_InHalfOpen_ReOpens(t *testing.T) {
	clk := &fixedClock{now: baseTime}
	b := newTestBreaker(1, 5*time.Second, clk)

	b.RecordFailure()
	clk.now = baseTime.Add(6 * time.Second)
	_ = b.Allow() // transitions to half-open

	b.RecordFailure()
	if b.CurrentState() != StateOpen {
		t.Fatalf("expected StateOpen after failure in half-open, got %v", b.CurrentState())
	}
}
