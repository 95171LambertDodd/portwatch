// Package sampler provides probabilistic sampling for port scan events,
// allowing high-volume environments to reduce alert noise by only forwarding
// a configurable fraction of events downstream.
package sampler

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/user/portwatch/internal/portscanner"
)

// Sampler decides whether a given port entry should be forwarded based on a
// sampling rate in the range (0.0, 1.0]. A rate of 1.0 forwards everything.
type Sampler struct {
	mu   sync.Mutex
	rng  *rand.Rand
	rate float64

	// stats
	total    int64
	forwarded int64
}

// New creates a Sampler with the given rate. Rate must be in (0.0, 1.0].
func New(rate float64) (*Sampler, error) {
	if rate <= 0.0 || rate > 1.0 {
		return nil, fmt.Errorf("sampler: rate must be in (0, 1], got %v", rate)
	}
	return &Sampler{
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
		rate: rate,
	}, nil
}

// Sample returns true if the entry should be forwarded downstream.
// It is safe for concurrent use.
func (s *Sampler) Sample(_ portscanner.PortEntry) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total++
	if s.rng.Float64() < s.rate {
		s.forwarded++
		return true
	}
	return false
}

// Stats returns the total number of entries seen and the number forwarded.
func (s *Sampler) Stats() (total, forwarded int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total, s.forwarded
}
