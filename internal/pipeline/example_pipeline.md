# Pipeline

The `pipeline` package wires all portwatch stages into a single `Run()` call.

## Usage

```go
p, err := pipeline.New(pipeline.Config{
    Scanner:  portscanner.NewScanner([]string{"/proc/net/tcp", "/proc/net/tcp6"}),
    Alerter:  alerting.NewAlerter(os.Stdout),
    Filter:   f,   // *filter.Filter — may be nil
    Dedup:    d,   // *dedup.Deduplicator — may be nil
    Suppress: s,   // *suppress.Suppressor — may be nil
    Notifier: n,   // *notify.Notifier — required
    Metrics:  m,   // *metrics.Metrics — may be nil
})
if err != nil {
    log.Fatal(err)
}

// Call Run on every watcher tick.
if err := p.Run(ctx); err != nil {
    log.Printf("pipeline error: %v", err)
}
```

## Stages

| Stage    | Required | Purpose                                      |
|----------|----------|----------------------------------------------|
| Scanner  | Yes      | Reads /proc/net/tcp* for active bindings     |
| Alerter  | No       | Diffs current vs previous snapshot          |
| Filter   | No       | Drops entries matching ignore rules         |
| Suppress | No       | Silences alerts covered by suppress rules   |
| Dedup    | No       | Prevents duplicate alerts within a TTL      |
| Notifier | Yes      | Delivers alerts to configured sinks         |
| Metrics  | No       | Records scan/alert counters                 |
