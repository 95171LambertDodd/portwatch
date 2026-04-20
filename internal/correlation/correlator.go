// Package correlation links port binding events to known services
// by matching port/protocol pairs against a configurable service map.
package correlation

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/portscanner"
)

// ServiceInfo holds metadata about a known service.
type ServiceInfo struct {
	Name        string
	Description string
	Expected    bool // true if this binding is expected in normal operation
}

// Rule maps a port+protocol pair to a ServiceInfo.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
	Service  ServiceInfo
}

// Correlator matches port entries against a set of known service rules.
type Correlator struct {
	rules map[string]ServiceInfo
}

// New creates a Correlator from the given rules.
// Returns an error if any rule has an invalid protocol.
func New(rules []Rule) (*Correlator, error) {
	m := make(map[string]ServiceInfo, len(rules))
	for _, r := range rules {
		proto := strings.ToLower(r.Protocol)
		if proto != "tcp" && proto != "udp" {
			return nil, fmt.Errorf("correlation: invalid protocol %q in rule for port %d", r.Protocol, r.Port)
		}
		key := ruleKey(r.Port, proto)
		m[key] = r.Service
	}
	return &Correlator{rules: m}, nil
}

// Lookup returns the ServiceInfo for the given port entry, and whether a match was found.
func (c *Correlator) Lookup(e portscanner.PortEntry) (ServiceInfo, bool) {
	key := ruleKey(e.Port, strings.ToLower(e.Protocol))
	svc, ok := c.rules[key]
	return svc, ok
}

// IsExpected returns true if the entry matches a rule marked as expected.
func (c *Correlator) IsExpected(e portscanner.PortEntry) bool {
	svc, ok := c.Lookup(e)
	return ok && svc.Expected
}

func ruleKey(port uint16, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
