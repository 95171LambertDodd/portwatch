package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	ScanCount      int64
	AlertCount     int64
	ViolationCount int64
	LastScanAt     time.Time
	LastAlertAt    time.Time
}

// Collector accumulates runtime metrics for portwatch.
type Collector struct {
	mu             sync.RWMutex
	scanCount      int64
	alertCount     int64
	violationCount int64
	lastScanAt     time.Time
	lastAlertAt    time.Time
	now            func() time.Time
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{now: time.Now}
}

// RecordScan increments the scan counter and records the timestamp.
func (c *Collector) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scanCount++
	c.lastScanAt = c.now()
}

// RecordAlert increments the alert counter and records the timestamp.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertCount++
	c.lastAlertAt = c.now()
}

// RecordViolation increments the baseline-violation counter.
func (c *Collector) RecordViolation() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.violationCount++
}

// Snapshot returns a consistent copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScanCount:      c.scanCount,
		AlertCount:     c.alertCount,
		ViolationCount: c.violationCount,
		LastScanAt:     c.lastScanAt,
		LastAlertAt:    c.lastAlertAt,
	}
}

// Reset zeroes all counters (useful in tests).
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scanCount = 0
	c.alertCount = 0
	c.violationCount = 0
	c.lastScanAt = time.Time{}
	c.lastAlertAt = time.Time{}
}
