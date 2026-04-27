// Package masking provides field-level redaction for port scan entries
// before they are emitted to sinks or logs.
package masking

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule describes a single masking rule.
type Rule struct {
	// Field is the entry field to redact: "pid", "process", "addr".
	Field string
	// Replacement is the string substituted for the original value.
	Replacement string
}

// Masker applies a set of redaction rules to port scan entries.
type Masker struct {
	rules []Rule
}

var validFields = map[string]struct{}{
	"pid":     {},
	"process": {},
	"addr":    {},
}

// New constructs a Masker from the given rules.
// Returns an error if any rule references an unknown field or has an empty replacement.
func New(rules []Rule) (*Masker, error) {
	for i, r := range rules {
		field := strings.ToLower(strings.TrimSpace(r.Field))
		if _, ok := validFields[field]; !ok {
			return nil, fmt.Errorf("masking: rule %d: unknown field %q", i, r.Field)
		}
		if r.Replacement == "" {
			return nil, fmt.Errorf("masking: rule %d: replacement must not be empty", i)
		}
		rules[i].Field = field
	}
	return &Masker{rules: rules}, nil
}

// Apply returns a copy of the entry with masked fields.
func (m *Masker) Apply(e portscanner.PortEntry) portscanner.PortEntry {
	out := e
	for _, r := range m.rules {
		switch r.Field {
		case "pid":
			out.PID = 0
		case "process":
			out.Process = r.Replacement
		case "addr":
			out.LocalAddr = r.Replacement
		}
	}
	return out
}

// ApplyAll applies masking to a slice of entries and returns a new slice.
func (m *Masker) ApplyAll(entries []portscanner.PortEntry) []portscanner.PortEntry {
	out := make([]portscanner.PortEntry, len(entries))
	for i, e := range entries {
		out[i] = m.Apply(e)
	}
	return out
}
