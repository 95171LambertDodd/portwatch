package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestNew_ReturnsNonNil(t *testing.T) {
	d := debounce.New(10*time.Millisecond, func(string) {})
	if d == nil {
		t.Fatal("expected non-nil debouncer")
	}
}

func TestTrigger_CallbackFiredAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var got []string

	d := debounce.New(30*time.Millisecond, func(key string) {
		mu.Lock()
		got = append(got, key)
		mu.Unlock()
	})

	d.Trigger("port:8080")
	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || got[0] != "port:8080" {
		t.Fatalf("expected callback with port:8080, got %v", got)
	}
}

func TestTrigger_ResetsPreventsEarlyFire(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(50*time.Millisecond, func(string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Trigger rapidly — callback should fire only once
	for i := 0; i < 5; i++ {
		d.Trigger("key")
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(80 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 callback, got %d", count)
	}
}

func TestCancel_PreventsFire(t *testing.T) {
	fired := false
	d := debounce.New(40*time.Millisecond, func(string) { fired = true })

	d.Trigger("k")
	d.Cancel("k")
	time.Sleep(70 * time.Millisecond)

	if fired {
		t.Fatal("callback should not have fired after cancel")
	}
}

func TestPending_ReflectsActiveTimers(t *testing.T) {
	d := debounce.New(100*time.Millisecond, func(string) {})

	if d.Pending() != 0 {
		t.Fatal("expected 0 pending initially")
	}
	d.Trigger("a")
	d.Trigger("b")
	if d.Pending() != 2 {
		t.Fatalf("expected 2 pending, got %d", d.Pending())
	}
	d.Cancel("a")
	if d.Pending() != 1 {
		t.Fatalf("expected 1 pending after cancel, got %d", d.Pending())
	}
}
