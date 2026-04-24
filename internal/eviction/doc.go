// Package eviction implements a thread-safe LRU cache with optional TTL-based
// expiry, suitable for bounding memory use in long-running portwatch daemons.
//
// # Overview
//
// portwatch accumulates per-port state across scan cycles. Without bounds,
// maps tracking process names, labels, and enrichment results can grow
// indefinitely. The eviction cache provides a single reusable primitive:
//
//	c, err := eviction.New(1024, 5*time.Minute)
//	c.Set("tcp:8080", resolvedInfo)
//	if v, ok := c.Get("tcp:8080"); ok {
//		// use cached value
//	}
//
// # Eviction Policy
//
// Items are evicted in LRU order when capacity is reached. If a non-zero TTL
// is configured, entries that have exceeded their lifetime are treated as
// misses and removed lazily on access.
package eviction
