// Package fingerprint provides per-entry identity hashing for the portwatch
// change-detection pipeline.
//
// # Overview
//
// The Tracker maintains a map of stable entry keys (protocol + address + port)
// to their last-seen Fingerprint. A Fingerprint is a SHA-256 hex digest that
// covers the full content of a PortEntry including the owning PID.
//
// # Usage
//
//	tr := fingerprint.New()
//
//	for _, e := range currentEntries {
//	    if tr.Changed(e) {
//	        // emit mutation event — new binding or PID change
//	    }
//	}
//
// Entries that disappear between scans should be passed to Remove so that
// if the same port is later re-bound it is correctly flagged as new.
//
// # Relationship to digest
//
// The digest package computes a single fingerprint over an entire slice of
// entries and is used to decide whether a full scan produced any change at
// all. This package operates at the individual-entry level and can detect
// which specific entries mutated (e.g. PID restart on a well-known port).
package fingerprint
