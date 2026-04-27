// Package routing provides rule-based routing of port scan entries to
// named output channels. Each rule matches on protocol, port range, or
// tag and directs matching entries to a labelled sink.
package routing

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule describes a single routing rule.
type Rule struct {
	// Name is a human-readable label for the rule.
	Name string
	// Protocol filters by "tcp" or "udp". Empty means any.
	Protocol string
	// MinPort and MaxPort define an inclusive port range. Zero values are ignored.
	MinPort uint16
	MaxPort uint16
	// Tag matches entries whose Tag field equals this value. Empty means any.
	Tag string
	// Destination is the named channel this rule routes to.
	Destination string
}

// Router routes port scan entries to named destinations based on ordered rules.
type Router struct {
	rules []Rule
}

// New validates and builds a Router from the provided rules.
// Rules are evaluated in order; the first match wins.
func New(rules []Rule) (*Router, error) {
	for i, r := range rules {
		if r.Name == "" {
			return nil, fmt.Errorf("rule[%d]: name must not be empty", i)
		}
		if r.Destination == "" {
			return nil, fmt.Errorf("rule[%d] %q: destination must not be empty", i, r.Name)
		}
		if r.Protocol != "" && r.Protocol != "tcp" && r.Protocol != "udp" {
			return nil, fmt.Errorf("rule[%d] %q: invalid protocol %q", i, r.Name, r.Protocol)
		}
		if r.MaxPort > 0 && r.MinPort > r.MaxPort {
			return nil, fmt.Errorf("rule[%d] %q: minPort %d > maxPort %d", i, r.Name, r.MinPort, r.MaxPort)
		}
	}
	if len(rules) == 0 {
		return nil, errors.New("at least one rule is required")
	}
	return &Router{rules: rules}, nil
}

// Route returns the destination name for the given entry.
// If no rule matches, an empty string is returned.
func (r *Router) Route(e portscanner.PortEntry) string {
	for _, rule := range r.rules {
		if rule.Protocol != "" && rule.Protocol != e.Protocol {
			continue
		}
		if rule.MinPort > 0 && e.Port < rule.MinPort {
			continue
		}
		if rule.MaxPort > 0 && e.Port > rule.MaxPort {
			continue
		}
		if rule.Tag != "" && rule.Tag != e.Tag {
			continue
		}
		return rule.Destination
	}
	return ""
}
