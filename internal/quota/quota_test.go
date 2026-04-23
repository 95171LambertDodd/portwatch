package quota

import (
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

func (f *fixedClock) now() time.Time { return f.t }
func (f *fixedClock) advance(d time.Duration) { f.t = f.t.Add(d) }

func newTestEnforcer(t *testing.T, window time.Duration, limit int, clk *fixedClock) *Enforcer {
	t.Helper()
	e, err := newWithClock(window, limit, clk.now)
	if err != nil {
		t.Fatalf("newWithClock: %v", err)
	}
	return e
}

func TestNew_ValidParams_ReturnsNonNil(t *testing.T) {
	e, err := New(time.Minute, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Enforcer")
	}
}

func TestNew_ZeroWindow_ReturnsError(t *testing.T) {
	_, err := New(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNew_ZeroLimit_ReturnsError(t *testing.T) {
	_, err := New(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 3, clk)
	if !e.Allow("key1") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_WithinLimit_AllAllowed(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 3, clk)
	for i := 0; i < 3; i++ {
		if !e.Allow("key1") {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit_Denied(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 2, clk)
	e.Allow("k")
	e.Allow("k")
	if e.Allow("k") {
		t.Fatal("third call should be denied after limit of 2")
	}
}

func TestAllow_WindowExpiry_ResetsCount(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 1, clk)
	e.Allow("k") // exhausts quota
	if e.Allow("k") {
		t.Fatal("should be denied within window")
	}
	clk.advance(2 * time.Minute)
	if !e.Allow("k") {
		t.Fatal("should be allowed after window expiry")
	}
}

func TestAllow_DifferentKeys_AreIndependent(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 1, clk)
	e.Allow("a")
	if !e.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestRemaining_FullQuotaForNewKey(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 5, clk)
	if got := e.Remaining("new"); got != 5 {
		t.Fatalf("expected 5 remaining, got %d", got)
	}
}

func TestRemaining_DecreasesAfterAllow(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 5, clk)
	e.Allow("k")
	e.Allow("k")
	if got := e.Remaining("k"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestRemaining_ZeroWhenExhausted(t *testing.T) {
	clk := &fixedClock{t: time.Now()}
	e := newTestEnforcer(t, time.Minute, 2, clk)
	e.Allow("k")
	e.Allow("k")
	if got := e.Remaining("k"); got != 0 {
		t.Fatalf("expected 0 remaining, got %d", got)
	}
}
