// Package tag provides port entry tagging based on configurable label rules.
package tag

import (
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule maps a match condition to a label.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp", empty means any
	Label    string
}

// Tagger assigns labels to port entries based on rules.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger from the provided rules.
func New(rules []Rule) (*Tagger, error) {
	for i, r := range rules {
		if r.Label == "" {
			return nil, fmt.Errorf("rule %d: label must not be empty", i)
		}
		if r.Protocol != "" && r.Protocol != "tcp" && r.Protocol != "udp" {
			return nil, fmt.Errorf("rule %d: invalid protocol %q", i, r.Protocol)
		}
	}
	return &Tagger{rules: rules}, nil
}

// Tag returns all labels that match the given entry.
func (t *Tagger) Tag(e portscanner.PortEntry) []string {
	var labels []string
	for _, r := range t.rules {
		if r.Port != 0 && r.Port != e.Port {
			continue
		}
		if r.Protocol != "" && r.Protocol != e.Protocol {
			continue
		}
		labels = append(labels, r.Label)
	}
	return labels
}
