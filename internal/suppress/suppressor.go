package suppress

import (
	"sync"
	"time"
)

// Rule defines a suppression rule for a specific port/protocol combination.
type Rule struct {
	Port     int
	Protocol string // "tcp" or "udp"
	Until    time.Time
}

// Suppressor tracks suppression rules to silence alerts for known ports.
type Suppressor struct {
	mu    sync.Mutex
	rules map[string]time.Time
	clock func() time.Time
}

// New returns a new Suppressor.
func New() *Suppressor {
	return &Suppressor{
		rules: make(map[string]time.Time),
		clock: time.Now,
	}
}

// Suppress adds a suppression rule for the given port/protocol until the given time.
func (s *Suppressor) Suppress(port int, protocol string, until time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := ruleKey(port, protocol)
	s.rules[key] = until
}

// IsSuppressed returns true if alerts for the given port/protocol should be silenced.
func (s *Suppressor) IsSuppressed(port int, protocol string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := ruleKey(port, protocol)
	until, ok := s.rules[key]
	if !ok {
		return false
	}
	if s.clock().After(until) {
		delete(s.rules, key)
		return false
	}
	return true
}

// Clear removes all suppression rules.
func (s *Suppressor) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = make(map[string]time.Time)
}

func ruleKey(port int, protocol string) string {
	return protocol + ":" + itoa(port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
