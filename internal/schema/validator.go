// Package schema provides validation for port entry configurations
// against a defined schema, ensuring fields are within acceptable bounds.
package schema

import (
	"errors"
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule defines constraints applied to a port entry during validation.
type Rule struct {
	MinPort  uint16
	MaxPort  uint16
	Protocols []string // allowed protocols, e.g. ["tcp", "udp"]
}

// Validator checks port entries against a set of rules.
type Validator struct {
	rules []Rule
}

// New creates a Validator with the given rules.
// Returns an error if any rule is malformed.
func New(rules []Rule) (*Validator, error) {
	for i, r := range rules {
		if r.MinPort > r.MaxPort {
			return nil, fmt.Errorf("rule %d: MinPort %d exceeds MaxPort %d", i, r.MinPort, r.MaxPort)
		}
		for _, p := range r.Protocols {
			if p != "tcp" && p != "udp" {
				return nil, fmt.Errorf("rule %d: unsupported protocol %q", i, p)
			}
		}
	}
	return &Validator{rules: rules}, nil
}

// Validate checks a single port entry against all rules.
// Returns a non-nil error describing the first violation found.
func (v *Validator) Validate(e portscanner.PortEntry) error {
	if len(v.rules) == 0 {
		return nil
	}
	for _, r := range v.rules {
		if err := r.check(e); err != nil {
			return err
		}
	}
	return nil
}

// ValidateAll validates a slice of entries and returns all violations.
func (v *Validator) ValidateAll(entries []portscanner.PortEntry) []error {
	var errs []error
	for _, e := range entries {
		if err := v.Validate(e); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (r Rule) check(e portscanner.PortEntry) error {
	if e.Port < r.MinPort || e.Port > r.MaxPort {
		return fmt.Errorf("port %d out of allowed range [%d, %d]", e.Port, r.MinPort, r.MaxPort)
	}
	if len(r.Protocols) > 0 {
		for _, p := range r.Protocols {
			if e.Protocol == p {
				return nil
			}
		}
		return errors.New(fmt.Sprintf("protocol %q not in allowed list %v", e.Protocol, r.Protocols))
	}
	return nil
}
