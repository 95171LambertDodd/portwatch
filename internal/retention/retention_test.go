package retention

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func makeEntry(port int) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: "tcp", State: "LISTEN"}
}

func TestNew_ZeroPolicy_ReturnsError(t *testing.T) {
	_, err := New(Policy{})
	if err == nil {
		t.Fatal("expected error for zero policy")
	}
}

func TestNew_MaxAgeOnly_ReturnsNonNil(t *testing.T) {
	m, err := New(Policy{MaxAge: time.Minute})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestRecord_And_Entries_RoundTrip(t *testing.T) {
	m, _ := newWithClock(Policy{MaxCount: 10}, fixedClock(epoch))
	m.Record(makeEntry(8080))
	m.Record(makeEntry(9090))
	entries := m.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestPrune_MaxCount_RemovesOldest(t *testing.T) {
	m, _ := newWithClock(Policy{MaxCount: 2}, fixedClock(epoch))
	m.Record(makeEntry(1))
	m.Record(makeEntry(2))
	m.Record(makeEntry(3))
	entries := m.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after pruning, got %d", len(entries))
	}
	if entries[0].Port != 2 || entries[1].Port != 3 {
		t.Errorf("unexpected entries after prune: %+v", entries)
	}
}

func TestPrune_MaxAge_RemovesStale(t *testing.T) {
	now := epoch
	clk := &now
	m, _ := newWithClock(Policy{MaxAge: time.Minute}, func() time.Time { return *clk })

	m.Record(makeEntry(1000))

	// Advance clock beyond MaxAge
	*clk = epoch.Add(2 * time.Minute)
	m.Record(makeEntry(2000))

	entries := m.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry after age prune, got %d", len(entries))
	}
	if entries[0].Port != 2000 {
		t.Errorf("expected port 2000, got %d", entries[0].Port)
	}
}

func TestPrune_BothPolicies_CountWins(t *testing.T) {
	now := epoch
	clk := &now
	m, _ := newWithClock(Policy{MaxAge: time.Hour, MaxCount: 1}, func() time.Time { return *clk })

	m.Record(makeEntry(80))
	m.Record(makeEntry(443))

	entries := m.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 443 {
		t.Errorf("expected port 443, got %d", entries[0].Port)
	}
}
