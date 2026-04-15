package filter

import (
	"net"

	"github.com/user/portwatch/internal/portscanner"
)

// Rule defines a single filter rule that can exclude port entries from alerting.
type Rule struct {
	// Port matches a specific port number. 0 means any port.
	Port uint16
	// Protocol matches "tcp", "udp", or empty for any.
	Protocol string
	// CIDR matches source/bind addresses within a subnet. Empty means any.
	CIDR string

	parsedNet *net.IPNet
}

// Filter holds a set of rules used to suppress known/expected port bindings.
type Filter struct {
	rules []Rule
}

// New constructs a Filter from the provided rules, pre-parsing any CIDR fields.
func New(rules []Rule) (*Filter, error) {
	parsed := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.CIDR != "" {
			_, ipNet, err := net.ParseCIDR(r.CIDR)
			if err != nil {
				return nil, err
			}
			r.parsedNet = ipNet
		}
		parsed = append(parsed, r)
	}
	return &Filter{rules: parsed}, nil
}

// Allow returns true when the entry should be forwarded for alerting,
// and false when it is suppressed by at least one rule.
func (f *Filter) Allow(entry portscanner.PortEntry) bool {
	for _, r := range f.rules {
		if r.matches(entry) {
			return false
		}
	}
	return true
}

// matches reports whether a single rule covers the given entry.
func (r *Rule) matches(entry portscanner.PortEntry) bool {
	if r.Port != 0 && r.Port != entry.Port {
		return false
	}
	if r.Protocol != "" && r.Protocol != entry.Protocol {
		return false
	}
	if r.parsedNet != nil {
		ip := net.ParseIP(entry.LocalAddr)
		if ip == nil || !r.parsedNet.Contains(ip) {
			return false
		}
	}
	return true
}
