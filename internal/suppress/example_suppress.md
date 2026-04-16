# Suppressor

The `suppress` package allows portwatch to silence alerts for specific port/protocol
combinations for a defined duration. This is useful when a known service temporarily
binds to an unexpected port during maintenance or deployment.

## Usage

```go
s := suppress.New()

// Suppress alerts for port 8080/tcp for 10 minutes
s.Suppress(8080, "tcp", time.Now().Add(10*time.Minute))

// Check before emitting an alert
if !s.IsSuppressed(entry.Port, entry.Protocol) {
    alerter.Emit(entry)
}

// Clear all rules (e.g., on config reload)
s.Clear()
```

## Rule Expiry

Rules are lazily expired: when `IsSuppressed` is called and the rule's deadline
has passed, the rule is automatically removed from memory.

## Thread Safety

All methods are safe for concurrent use.
