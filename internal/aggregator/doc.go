// Package aggregator provides time-bucketed accumulation of port scan results.
//
// Scan entries are collected into fixed-duration Buckets. When a bucket's
// time window expires, it is promoted to a completed list and a new bucket
// is opened. Completed buckets can be flushed and forwarded to downstream
// sinks (e.g. JSON log files, metrics systems) via FlushSink.
//
// Typical usage:
//
//	agg, err := aggregator.New(30*time.Second, 10)
//	if err != nil { ... }
//
//	// on each scan tick:
//	agg.Add(entries)
//
//	// periodically:
//	for _, bucket := range agg.Flush() {
//		fmt.Println(bucket.Count, "entries in", bucket.Start)
//	}
package aggregator
