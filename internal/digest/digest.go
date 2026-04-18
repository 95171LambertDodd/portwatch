// Package digest computes and compares fingerprints of port scan snapshots
// to detect meaningful changes between polling cycles.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/user/portwatch/internal/portscanner"
)

// Fingerprint is a stable hash of a set of port entries.
type Fingerprint string

// Compute returns a deterministic SHA-256 fingerprint for the given entries.
// Entries are sorted before hashing so order does not matter.
func Compute(entries []portscanner.Entry) (Fingerprint, error) {
	sorted := make([]portscanner.Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Protocol != sorted[j].Protocol {
			return sorted[i].Protocol < sorted[j].Protocol
		}
		if sorted[i].LocalAddr != sorted[j].LocalAddr {
			return sorted[i].LocalAddr < sorted[j].LocalAddr
		}
		return sorted[i].LocalPort < sorted[j].LocalPort
	})

	b, err := json.Marshal(sorted)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(b)
	return Fingerprint(hex.EncodeToString(sum[:])), nil
}

// Equal returns true when two fingerprints match.
func Equal(a, b Fingerprint) bool {
	return a == b
}
