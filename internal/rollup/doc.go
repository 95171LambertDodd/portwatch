// Package rollup provides alert aggregation for portwatch.
//
// When the same port binding triggers many alerts in a short window,
// rollup groups them into a single summarized event to reduce noise.
//
// Usage:
//
//	agg := rollup.New(30 * time.Second)
//
//	// On each alert:
//	agg.Record("tcp:8080")
//
//	// Periodically (or via FlushSink):
//	events := agg.Flush()
//	for _, e := range events {
//	    fmt.Println(rollup.Summary(e))
//	}
package rollup
