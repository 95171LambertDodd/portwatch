package schema_test

import (
	"testing"

	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/schema"
)

func makeEntry(port uint16, proto string) portscanner.PortEntry {
	return portscanner.PortEntry{Port: port, Protocol: proto}
}

func TestNew_ValidRules_ReturnsNonNil(t *testing.T) {
	v, err := schema.New([]schema.Rule{
		{MinPort: 1024, MaxPort: 65535, Protocols: []string{"tcp"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
}

func TestNew_InvalidPortRange_ReturnsError(t *testing.T) {
	_, err := schema.New([]schema.Rule{
		{MinPort: 9000, MaxPort: 1000},
	})
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}

func TestNew_InvalidProtocol_ReturnsError(t *testing.T) {
	_, err := schema.New([]schema.Rule{
		{MinPort: 1, MaxPort: 65535, Protocols: []string{"sctp"}},
	})
	if err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestValidate_PortInRange_NoError(t *testing.T) {
	v, _ := schema.New([]schema.Rule{
		{MinPort: 1024, MaxPort: 9000, Protocols: []string{"tcp", "udp"}},
	})
	if err := v.Validate(makeEntry(8080, "tcp")); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_PortOutOfRange_ReturnsError(t *testing.T) {
	v, _ := schema.New([]schema.Rule{
		{MinPort: 1024, MaxPort: 9000, Protocols: []string{"tcp"}},
	})
	if err := v.Validate(makeEntry(80, "tcp")); err == nil {
		t.Error("expected error for port out of range")
	}
}

func TestValidate_DisallowedProtocol_ReturnsError(t *testing.T) {
	v, _ := schema.New([]schema.Rule{
		{MinPort: 1, MaxPort: 65535, Protocols: []string{"tcp"}},
	})
	if err := v.Validate(makeEntry(53, "udp")); err == nil {
		t.Error("expected error for disallowed protocol")
	}
}

func TestValidateAll_ReturnsAllViolations(t *testing.T) {
	v, _ := schema.New([]schema.Rule{
		{MinPort: 1024, MaxPort: 65535, Protocols: []string{"tcp"}},
	})
	entries := []portscanner.PortEntry{
		makeEntry(80, "tcp"),  // port violation
		makeEntry(8080, "tcp"), // ok
		makeEntry(443, "tcp"),  // port violation
	}
	errs := v.ValidateAll(entries)
	if len(errs) != 2 {
		t.Errorf("expected 2 violations, got %d", len(errs))
	}
}

func TestValidate_NoRules_AlwaysPasses(t *testing.T) {
	v, _ := schema.New(nil)
	if err := v.Validate(makeEntry(1, "udp")); err != nil {
		t.Errorf("expected no error with empty rules, got: %v", err)
	}
}
