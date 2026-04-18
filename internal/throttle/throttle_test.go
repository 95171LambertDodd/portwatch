package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) throttle.Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	th := throttle.New(time.Minute, 3, fixedClock(epoch))
	if !th.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_BurstRespected(t *testing.T) {
	th := throttle.New(time.Minute, 3, fixedClock(epoch))
	for i := 0; i < 3; i++ {
		if !th.Allow("k") {
			t.Fatalf("expected call %d to be allowed", i+1)
		}
	}
	if th.Allow("k") {
		t.Fatal("expected call beyond burst to be denied")
	}
}

func TestAllow_WindowResetAllows(t *testing.T) {
	var now = epoch
	clock := func() time.Time { return now }
	th := throttle.New(time.Minute, 2, clock)
	th.Allow("k")
	th.Allow("k")
	if th.Allow("k") {
		t.Fatal("should be denied within window")
	}
	now = epoch.Add(2 * time.Minute)
	if !th.Allow("k") {
		t.Fatal("should be allowed after window reset")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(time.Minute, 1, fixedClock(epoch))
	th.Allow("a")
	if !th.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := throttle.New(time.Minute, 1, fixedClock(epoch))
	th.Allow("k")
	if th.Allow("k") {
		t.Fatal("should be denied after burst")
	}
	th.Reset("k")
	if !th.Allow("k") {
		t.Fatal("should be allowed after reset")
	}
}

func TestStats_ReturnsCountAndSince(t *testing.T) {
	th := throttle.New(time.Minute, 5, fixedClock(epoch))
	th.Allow("k")
	th.Allow("k")
	count, since := th.Stats("k")
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
	if !since.Equal(epoch) {
		t.Fatalf("expected since %v, got %v", epoch, since)
	}
}

func TestStats_UnknownKey(t *testing.T) {
	th := throttle.New(time.Minute, 5, fixedClock(epoch))
	count, since := th.Stats("missing")
	if count != 0 || !since.IsZero() {
		t.Fatal("expected zero values for unknown key")
	}
}
