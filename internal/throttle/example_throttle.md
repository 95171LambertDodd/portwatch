# Throttle

The `throttle` package provides a burst-aware rate limiter keyed by arbitrary strings (e.g. port+protocol combos).

## Concepts

- **Window**: the time duration over which burst count is tracked.
- **MaxBurst**: maximum number of allowed actions within a single window.
- Once the burst is exhausted the key is blocked until the window expires.

## Example

```go
th := throttle.New(30*time.Second, 3, nil)

key := "tcp:8080"
if th.Allow(key) {
    // emit alert
}
```

## Difference from ratelimit.Limiter

| Feature | ratelimit | throttle |
|---------|-----------|----------|
| Cooldown (one-at-a-time) | ✓ | – |
| Burst window | – | ✓ |

Use `throttle` when you want to allow a small burst of alerts before silencing.
