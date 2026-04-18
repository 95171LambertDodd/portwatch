// Package throttle provides a burst-aware, key-scoped action throttler
// for use in portwatch alerting pipelines.
//
// Unlike a simple cooldown limiter, throttle tracks a count of actions
// within a sliding window and allows up to MaxBurst occurrences before
// blocking further actions until the window expires.
//
// It is safe for concurrent use.
package throttle
