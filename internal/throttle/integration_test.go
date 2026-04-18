package throttle_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

// TestAllow_MultipleKeys_BurstIndependent verifies that burst counters are
// tracked independently across many keys simultaneously.
func TestAllow_MultipleKeys_BurstIndependent(t *testing.T) {
	th := throttle.New(time.Minute, 2, fixedClock(epoch))
	keys := []string{"tcp:80", "tcp:443", "udp:53"}

	for _, k := range keys {
		for i := 0; i < 2; i++ {
			if !th.Allow(k) {
				t.Fatalf("key %s call %d should be allowed", k, i+1)
			}
		}
		if th.Allow(k) {
			t.Fatalf("key %s should be denied after burst", k)
		}
	}
}

// TestAllow_BurstThenReset_ThenBurstAgain simulates a full cycle.
func TestAllow_BurstThenReset_ThenBurstAgain(t *testing.T) {
	var now = epoch
	clock := func() time.Time { return now }
	th := throttle.New(10*time.Second, 2, clock)

	results := []bool{}
	for i := 0; i < 4; i++ {
		results = append(results, th.Allow("k"))
	}
	// advance past window
	now = epoch.Add(15 * time.Second)
	for i := 0; i < 2; i++ {
		results = append(results, th.Allow("k"))
	}

	expected := []bool{true, true, false, false, true, true}
	for i, r := range results {
		if r != expected[i] {
			t.Errorf("call %d: got %v want %v", i, r, expected[i])
		}
	}
}

// TestStats_TracksCorrectly verifies stats after several Allow calls.
func TestStats_TracksCorrectly(t *testing.T) {
	th := throttle.New(time.Minute, 10, fixedClock(epoch))
	key := fmt.Sprintf("udp:%d", 5353)
	for i := 0; i < 5; i++ {
		th.Allow(key)
	}
	count, since := th.Stats(key)
	if count != 5 {
		t.Fatalf("expected 5 got %d", count)
	}
	if since != epoch {
		t.Fatalf("unexpected since: %v", since)
	}
}
