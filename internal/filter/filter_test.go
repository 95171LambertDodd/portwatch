package filter

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
)

func entry(proto, addr string, port uint16) portscanner.PortEntry {
	return portscanner.PortEntry{
		Protocol:  proto,
		LocalAddr: addr,
		Port:      port,
	}
}

func TestNew_NoRules_AllowsEverything(t *testing.T) {
	f, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(entry("tcp", "127.0.0.1", 8080)) {
		t.Error("expected entry to be allowed with no rules")
	}
}

func TestNew_InvalidCIDR_ReturnsError(t *testing.T) {
	_, err := New([]Rule{{CIDR: "not-a-cidr"}})
	if err == nil {
		t.Fatal("expected error for invalid CIDR, got nil")
	}
}

func TestAllow_MatchesExactPort(t *testing.T) {
	f, _ := New([]Rule{{Port: 22}})
	if f.Allow(entry("tcp", "0.0.0.0", 22)) {
		t.Error("expected port 22 to be suppressed")
	}
	if !f.Allow(entry("tcp", "0.0.0.0", 80)) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestAllow_MatchesProtocol(t *testing.T) {
	f, _ := New([]Rule{{Protocol: "udp"}})
	if f.Allow(entry("udp", "0.0.0.0", 53)) {
		t.Error("expected udp entry to be suppressed")
	}
	if !f.Allow(entry("tcp", "0.0.0.0", 53)) {
		t.Error("expected tcp entry to be allowed")
	}
}

func TestAllow_MatchesCIDR(t *testing.T) {
	f, _ := New([]Rule{{CIDR: "127.0.0.0/8"}})
	if f.Allow(entry("tcp", "127.0.0.1", 9000)) {
		t.Error("expected loopback address to be suppressed")
	}
	if !f.Allow(entry("tcp", "192.168.1.5", 9000)) {
		t.Error("expected non-loopback address to be allowed")
	}
}

func TestAllow_CombinedRule_AllFieldsMustMatch(t *testing.T) {
	f, _ := New([]Rule{{Port: 443, Protocol: "tcp", CIDR: "0.0.0.0/0"}})
	// Port matches but protocol differs — should be allowed.
	if !f.Allow(entry("udp", "0.0.0.0", 443)) {
		t.Error("expected udp/443 to be allowed (protocol mismatch)")
	}
	// All fields match — should be suppressed.
	if f.Allow(entry("tcp", "10.0.0.1", 443)) {
		t.Error("expected tcp/443 in 0.0.0.0/0 to be suppressed")
	}
}

func TestAllow_MultipleRules_AnyMatchSuppresses(t *testing.T) {
	f, _ := New([]Rule{
		{Port: 22},
		{Port: 3306},
	})
	if f.Allow(entry("tcp", "0.0.0.0", 3306)) {
		t.Error("expected port 3306 to be suppressed")
	}
	if !f.Allow(entry("tcp", "0.0.0.0", 8080)) {
		t.Error("expected port 8080 to be allowed")
	}
}
