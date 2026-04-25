// Package fingerprint provides per-entry identity hashing for change detection.
// It assigns a stable, content-derived key to each port entry so that the
// watcher can detect mutations (e.g. PID change on same port) independently
// of the digest package which operates on slices.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/user/portwatch/internal/portscanner"
)

// Fingerprint is a stable hex string derived from a single PortEntry.
type Fingerprint string

// Tracker maintains the last-seen fingerprint for each port entry key.
type Tracker struct {
	seen map[string]Fingerprint
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{seen: make(map[string]Fingerprint)}
}

// Compute returns a Fingerprint for the given entry.
// The fingerprint covers: protocol, address, port, and PID so that a
// process restart on the same port is detected as a change.
func Compute(e portscanner.PortEntry) Fingerprint {
	raw := fmt.Sprintf("%s|%s|%d|%d", e.Protocol, e.LocalAddr, e.LocalPort, e.PID)
	sum := sha256.Sum256([]byte(raw))
	return Fingerprint(hex.EncodeToString(sum[:]))
}

// entryKey returns the map key used to track a port entry across scans.
func entryKey(e portscanner.PortEntry) string {
	return fmt.Sprintf("%s:%s:%d", e.Protocol, e.LocalAddr, e.LocalPort)
}

// Changed reports whether the entry's fingerprint differs from the last
// recorded value. It always updates the internal state to the new fingerprint.
func (t *Tracker) Changed(e portscanner.PortEntry) bool {
	key := entryKey(e)
	next := Compute(e)
	prev, ok := t.seen[key]
	t.seen[key] = next
	if !ok {
		// First time we have seen this entry — treat as changed so callers
		// can emit a "new binding" event.
		return true
	}
	return prev != next
}

// Remove deletes the tracked fingerprint for an entry. Call this when an
// entry disappears so that if it reappears it is treated as new.
func (t *Tracker) Remove(e portscanner.PortEntry) {
	delete(t.seen, entryKey(e))
}

// Len returns the number of entries currently tracked.
func (t *Tracker) Len() int {
	return len(t.seen)
}
