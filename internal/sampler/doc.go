// Package sampler implements probabilistic sampling for portwatch scan events.
//
// # Overview
//
// In high-traffic environments the port scanner may detect hundreds of entries
// per tick. The Sampler allows operators to forward only a representative
// fraction of those entries to downstream consumers (notifiers, history
// recorders, exporters) without losing statistical coverage.
//
// # Usage
//
//	s, err := sampler.New(0.25) // forward ~25 % of entries
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range entries {
//		if s.Sample(entry) {
//			notifier.Notify(entry)
//		}
//	}
//	total, forwarded := s.Stats()
//	fmt.Printf("sampled %d/%d entries\n", forwarded, total)
//
// # Thread Safety
//
// Sampler is safe for concurrent use by multiple goroutines.
package sampler
