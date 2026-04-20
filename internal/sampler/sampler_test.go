package sampler

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func makeEntry(port uint16) portscanner.PortEntry {
	return portscanner.PortEntry{
		Protocol:  "tcp",
		LocalPort: port,
		State:     "LISTEN",
	}
}

func TestNew_ValidRate_ReturnsNonNil(t *testing.T) {
	s, err := New(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestNew_ZeroRate_ReturnsError(t *testing.T) {
	_, err := New(0.0)
	if err == nil {
		t.Fatal("expected error for rate=0.0")
	}
}

func TestNew_NegativeRate_ReturnsError(t *testing.T) {
	_, err := New(-0.1)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNew_RateAboveOne_ReturnsError(t *testing.T) {
	_, err := New(1.1)
	if err == nil {
		t.Fatal("expected error for rate > 1.0")
	}
}

func TestSample_RateOne_AlwaysForwards(t *testing.T) {
	s, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 100; i++ {
		if !s.Sample(makeEntry(uint16(i % 65535 + 1))) {
			t.Fatalf("expected all entries to be forwarded at rate=1.0")
		}
	}
}

func TestSample_StatsTracked(t *testing.T) {
	s, err := New(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 10; i++ {
		s.Sample(makeEntry(8080))
	}
	total, forwarded := s.Stats()
	if total != 10 {
		t.Errorf("expected total=10, got %d", total)
	}
	if forwarded != 10 {
		t.Errorf("expected forwarded=10, got %d", forwarded)
	}
}

func TestSample_LowRate_SomeDrop(t *testing.T) {
	s, err := New(0.01)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 1000; i++ {
		s.Sample(makeEntry(uint16(i%1000 + 1)))
	}
	total, forwarded := s.Stats()
	if total != 1000 {
		t.Errorf("expected total=1000, got %d", total)
	}
	// With rate=0.01 over 1000 samples we should drop the vast majority.
	if forwarded >= 100 {
		t.Errorf("expected far fewer than 100 forwarded at rate=0.01, got %d", forwarded)
	}
}
