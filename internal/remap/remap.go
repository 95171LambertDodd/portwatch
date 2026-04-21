// Package remap provides port remapping support, allowing known ports to be
// aliased to canonical names or forwarded to alternate addresses for alerting.
package remap

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule defines a single remapping rule.
type Rule struct {
	// FromPort is the original port to match.
	FromPort uint16
	// Protocol must be "tcp" or "udp".
	Protocol string
	// Alias is the canonical name assigned to this port.
	Alias string
	// ToPort is an optional alternate port to remap to (0 means no remap).
	ToPort uint16
}

// Result holds the outcome of a remap lookup.
type Result struct {
	Alias  string
	ToPort uint16
	Mapped bool
}

// Remapper applies port remapping rules to scanner entries.
type Remapper struct {
	rules map[ruleKey]Rule
}

type ruleKey struct {
	port  uint16
	proto string
}

// New creates a Remapper from the given rules.
// Returns an error if any rule has an invalid protocol or empty alias.
func New(rules []Rule) (*Remapper, error) {
	if len(rules) == 0 {
		return &Remapper{rules: make(map[ruleKey]Rule)}, nil
	}

	m := make(map[ruleKey]Rule, len(rules))
	for i, r := range rules {
		if r.Protocol != "tcp" && r.Protocol != "udp" {
			return nil, fmt.Errorf("rule %d: invalid protocol %q, must be tcp or udp", i, r.Protocol)
		}
		if r.Alias == "" {
			return nil, errors.New(fmt.Sprintf("rule %d: alias must not be empty", i))
		}
		key := ruleKey{port: r.FromPort, proto: r.Protocol}
		m[key] = r
	}
	return &Remapper{rules: m}, nil
}

// Lookup returns the remap Result for a given scanner entry.
// If no rule matches, Result.Mapped is false.
func (r *Remapper) Lookup(e portscanner.PortEntry) Result {
	key := ruleKey{port: uint16(e.Port), proto: e.Protocol}
	rule, ok := r.rules[key]
	if !ok {
		return Result{}
	}
	return Result{
		Alias:  rule.Alias,
		ToPort: rule.ToPort,
		Mapped: true,
	}
}

// Apply returns a copy of the entry with Port replaced by ToPort if a remap
// rule with a non-zero ToPort exists. The original entry is unchanged.
func (r *Remapper) Apply(e portscanner.PortEntry) portscanner.PortEntry {
	res := r.Lookup(e)
	if res.Mapped && res.ToPort != 0 {
		e.Port = int(res.ToPort)
	}
	return e
}
