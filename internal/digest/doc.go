// Package digest provides deterministic fingerprinting of port scan results.
//
// A Fingerprint is a SHA-256 hash of a canonicalised (sorted) list of
// portscanner.Entry values serialised as JSON. Two scans that produce
// identical bindings — regardless of the order in which the OS reports
// them — will yield the same Fingerprint, allowing the watcher to skip
// downstream processing when nothing has changed.
//
// Typical usage:
//
//	prev, _ := digest.Compute(lastEntries)
//	curr, _ := digest.Compute(newEntries)
//	if digest.Equal(prev, curr) {
//	    return // no change
//	}
package digest
