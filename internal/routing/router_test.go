package routing_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/routing"
)

func entry(proto string, port uint16, tag string) portscanner.PortEntry {
	return portscanner.PortEntry{Protocol: proto, Port: port, Tag: tag}
}

func TestNew_ValidRules_ReturnsNonNil(t *testing.T) {
	r, err := routing.New([]routing.Rule{
		{Name: "all", Destination: "default"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestNew_EmptyRules_ReturnsError(t *testing.T) {
	_, err := routing.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_MissingName_ReturnsError(t *testing.T) {
	_, err := routing.New([]routing.Rule{{Destination: "sink"}})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestNew_MissingDestination_ReturnsError(t *testing.T) {
	_, err := routing.New([]routing.Rule{{Name: "r1"}})
	if err == nil {
		t.Fatal("expected error for missing destination")
	}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := routing.New([]routing.Rule{
		{Name: "r1", Protocol: "quic", Destination: "sink"},
	})
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestNew_InvertedPortRange_ReturnsError(t *testing.T) {
	_, err := routing.New([]routing.Rule{
		{Name: "r1", MinPort: 9000, MaxPort: 80, Destination: "sink"},
	})
	if err == nil {
		t.Fatal("expected error for inverted port range")
	}
}

func TestRoute_MatchesProtocol(t *testing.T) {
	r, _ := routing.New([]routing.Rule{
		{Name: "tcp-only", Protocol: "tcp", Destination: "tcp-sink"},
		{Name: "fallback", Destination: "default"},
	})
	if got := r.Route(entry("tcp", 80, "")); got != "tcp-sink" {
		t.Errorf("expected tcp-sink, got %q", got)
	}
	if got := r.Route(entry("udp", 80, "")); got != "default" {
		t.Errorf("expected default, got %q", got)
	}
}

func TestRoute_MatchesPortRange(t *testing.T) {
	r, _ := routing.New([]routing.Rule{
		{Name: "high", MinPort: 1024, MaxPort: 65535, Destination: "high-sink"},
		{Name: "low", Destination: "low-sink"},
	})
	if got := r.Route(entry("tcp", 8080, "")); got != "high-sink" {
		t.Errorf("expected high-sink, got %q", got)
	}
	if got := r.Route(entry("tcp", 22, "")); got != "low-sink" {
		t.Errorf("expected low-sink, got %q", got)
	}
}

func TestRoute_MatchesTag(t *testing.T) {
	r, _ := routing.New([]routing.Rule{
		{Name: "critical", Tag: "critical", Destination: "alert-sink"},
		{Name: "fallback", Destination: "default"},
	})
	if got := r.Route(entry("tcp", 443, "critical")); got != "alert-sink" {
		t.Errorf("expected alert-sink, got %q", got)
	}
	if got := r.Route(entry("tcp", 443, "info")); got != "default" {
		t.Errorf("expected default, got %q", got)
	}
}

func TestRoute_NoMatchReturnsEmpty(t *testing.T) {
	r, _ := routing.New([]routing.Rule{
		{Name: "tcp-only", Protocol: "tcp", Destination: "tcp-sink"},
	})
	if got := r.Route(entry("udp", 53, "")); got != "" {
		t.Errorf("expected empty destination, got %q", got)
	}
}
