// Package circuitbreaker provides a thread-safe circuit breaker for use in
// portwatch pipelines.
//
// A Breaker starts in the Closed state, allowing all calls through. After a
// configurable number of consecutive failures it transitions to Open, blocking
// further calls and returning ErrOpen. Once the cooldown period has elapsed
// the breaker moves to HalfOpen, permitting a single probe call. A successful
// probe resets the breaker to Closed; another failure returns it to Open.
//
// Typical usage:
//
//	b := circuitbreaker.New(5, 30*time.Second)
//
//	if err := b.Allow(); err != nil {
//		// skip processing — circuit is open
//		return
//	}
//	if err := doWork(); err != nil {
//		b.RecordFailure()
//		return err
//	}
//	b.RecordSuccess()
package circuitbreaker
