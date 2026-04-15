package ratelimit

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallAlwaysAllowed(t *testing.T) {
	l := New(5 * time.Second)
	if !l.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownDenied(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("key1")

	// Advance time by less than cooldown
	l.now = fixedClock(base.Add(5 * time.Second))
	if l.Allow("key1") {
		t.Fatal("expected second call within cooldown to be denied")
	}
}

func TestAllow_CallAfterCooldownAllowed(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("key1")

	// Advance time beyond cooldown
	l.now = fixedClock(base.Add(11 * time.Second))
	if !l.Allow("key1") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("key1")

	l.now = fixedClock(base.Add(1 * time.Second))
	if !l.Allow("key2") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsKeyImmediately(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("key1")
	l.Reset("key1")

	// Same time, but reset should allow it
	if !l.Allow("key1") {
		t.Fatal("expected key to be allowed after reset")
	}
}

func TestFlush_ClearsAllKeys(t *testing.T) {
	base := time.Now()
	l := New(10 * time.Second)
	l.now = fixedClock(base)

	l.Allow("key1")
	l.Allow("key2")
	l.Flush()

	if !l.Allow("key1") {
		t.Fatal("expected key1 to be allowed after flush")
	}
	if !l.Allow("key2") {
		t.Fatal("expected key2 to be allowed after flush")
	}
}
