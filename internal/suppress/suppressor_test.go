package suppress

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestIsSuppressed_NotSuppressed(t *testing.T) {
	s := New()
	if s.IsSuppressed(8080, "tcp") {
		t.Fatal("expected not suppressed")
	}
}

func TestSuppress_ActiveRule(t *testing.T) {
	s := New()
	now := time.Now()
	s.clock = fixedClock(now)
	s.Suppress(8080, "tcp", now.Add(1*time.Minute))
	if !s.IsSuppressed(8080, "tcp") {
		t.Fatal("expected suppressed")
	}
}

func TestSuppress_ExpiredRule(t *testing.T) {
	s := New()
	now := time.Now()
	s.Suppress(8080, "tcp", now.Add(-1*time.Second))
	s.clock = fixedClock(now)
	if s.IsSuppressed(8080, "tcp") {
		t.Fatal("expected expired rule to not suppress")
	}
}

func TestSuppress_DifferentProtocols(t *testing.T) {
	s := New()
	now := time.Now()
	s.clock = fixedClock(now)
	s.Suppress(53, "udp", now.Add(1*time.Minute))
	if s.IsSuppressed(53, "tcp") {
		t.Fatal("tcp should not be suppressed")
	}
	if !s.IsSuppressed(53, "udp") {
		t.Fatal("udp should be suppressed")
	}
}

func TestClear_RemovesAllRules(t *testing.T) {
	s := New()
	now := time.Now()
	s.clock = fixedClock(now)
	s.Suppress(80, "tcp", now.Add(1*time.Hour))
	s.Suppress(443, "tcp", now.Add(1*time.Hour))
	s.Clear()
	if s.IsSuppressed(80, "tcp") || s.IsSuppressed(443, "tcp") {
		t.Fatal("expected all rules cleared")
	}
}

func TestExpiredRule_DeletedFromMap(t *testing.T) {
	s := New()
	now := time.Now()
	s.Suppress(9000, "tcp", now.Add(-1*time.Second))
	s.clock = fixedClock(now)
	s.IsSuppressed(9000, "tcp") // trigger cleanup
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[ruleKey(9000, "tcp")]; ok {
		t.Fatal("expected expired rule to be deleted")
	}
}
