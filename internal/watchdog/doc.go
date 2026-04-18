// Package watchdog provides a heartbeat monitor for portwatch scan cycles.
//
// A Watchdog is created with a maximum allowed age. The watcher loop calls
// Beat after each successful scan. Check can be polled (e.g. by the health
// endpoint) to determine whether scans are running on schedule.
//
// Example usage:
//
//	wd := watchdog.New(30 * time.Second)
//
//	// inside scan loop:
//	wd.Beat()
//
//	// inside health check:
//	status, msg := wd.Check()
package watchdog
