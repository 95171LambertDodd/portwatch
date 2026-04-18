// Package trend provides a sliding-window frequency tracker for port binding events.
//
// Use [New] to create a Tracker with a retention window. Call [Tracker.Record] each
// scan cycle with the number of active bindings for a key (e.g. "tcp:8080").
// Call [Tracker.Rate] to retrieve the average count observed within the window,
// which can be compared against a threshold to detect sudden spikes.
//
// Example:
//
//	tr := trend.New(5 * time.Minute)
//	tr.Record("tcp:8080", len(entries))
//	if tr.Rate("tcp:8080") > float64(cfg.SpikeThreshold) {
//		// emit alert
//	}
package trend
