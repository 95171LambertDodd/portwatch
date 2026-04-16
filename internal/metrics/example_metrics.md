# Metrics

The `metrics` package provides a thread-safe `Collector` that tracks runtime
counters for the portwatch daemon.

## Counters

| Counter        | Incremented by          |
|----------------|-------------------------|
| `ScanCount`    | Each port-scan cycle    |
| `AlertCount`   | Each alert emitted      |
| `ViolationCount` | Each baseline violation |

## Usage

```go
col := metrics.New()

// inside the scan loop
col.RecordScan()

// when an alert is emitted
col.RecordAlert()

// when a baseline violation is detected
col.RecordViolation()

// read current values (safe for concurrent access)
snap := col.Snapshot()
fmt.Printf("scans=%d alerts=%d violations=%d\n",
    snap.ScanCount, snap.AlertCount, snap.ViolationCount)
```

## Notes

- All methods are safe for concurrent use.
- `Reset()` is provided primarily for use in tests.
- Timestamps (`LastScanAt`, `LastAlertAt`) are zero-valued until the
  corresponding event is recorded.
