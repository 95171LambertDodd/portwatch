package scorecard

import (
	"github.com/user/portwatch/internal/portscanner"
)

// DefaultRules returns a sensible set of built-in scoring rules suitable for
// most deployments. Callers may append additional rules before constructing a
// Scorer.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:   "privileged-port",
			Weight: 30,
			Match: func(e portscanner.Entry) bool {
				return e.Port > 0 && e.Port < 1024
			},
		},
		{
			Name:   "high-ephemeral-port",
			Weight: 5,
			Match: func(e portscanner.Entry) bool {
				return e.Port >= 49152
			},
		},
		{
			Name:   "udp-binding",
			Weight: 15,
			Match: func(e portscanner.Entry) bool {
				return e.Protocol == "udp"
			},
		},
		{
			Name:   "pid-zero",
			Weight: 50,
			Match: func(e portscanner.Entry) bool {
				return e.PID == 0
			},
		},
	}
}
