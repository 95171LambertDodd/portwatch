# eviction — LRU Cache with TTL

The `eviction` package provides a bounded, thread-safe LRU cache used by
portwatch internals to limit memory growth across long-running scan cycles.

## Usage

```go
import "github.com/yourorg/portwatch/internal/eviction"

// Create a cache of up to 512 entries, each valid for 2 minutes.
cache, err := eviction.New(512, 2*time.Minute)
if err != nil {
    log.Fatal(err)
}

// Store a resolved process name for a port key.
cache.Set("tcp:9090", "prometheus")

// Retrieve it later.
if name, ok := cache.Get("tcp:9090"); ok {
    fmt.Println("cached:", name)
}
```

## Eviction Behaviour

| Condition | Result |
|-----------|--------|
| Capacity reached | Least-recently-used entry removed |
| TTL exceeded | Entry treated as miss, removed on access |
| Zero TTL | Time-based expiry disabled |
| Key updated | LRU order refreshed, size unchanged |

## Integration

The cache is used by `internal/process` to avoid repeated `/proc` lookups
and by `internal/enrichment` to cache reverse-DNS results between scan ticks.
