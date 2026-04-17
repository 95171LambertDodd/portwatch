# Audit Log

The `audit` package provides append-only structured logging of significant portwatch events.

## Event Kinds

| Kind | Description |
|------|-------------|
| `new_binding` | A port was bound that was not previously seen |
| `violation` | A port binding violates the baseline policy |
| `suppressed` | An alert was suppressed by a suppression rule |
| `scan_error` | An error occurred during a port scan |

## File Format

Events are stored as newline-delimited JSON (JSONL):

```json
{"timestamp":"2024-01-01T00:00:00Z","kind":"new_binding","port":8080,"protocol":"tcp","pid":1234,"message":"unexpected binding detected"}
{"timestamp":"2024-01-01T00:01:00Z","kind":"violation","port":443,"protocol":"tcp","pid":999,"message":"port not in baseline"}
```

## Usage

```go
logger, err := audit.NewLogger("/var/log/portwatch/audit.log")
if err != nil {
    log.Fatal(err)
}

logger.Log(audit.Event{
    Kind:     "new_binding",
    Port:     8080,
    Protocol: "tcp",
    PID:      1234,
    Message:  "unexpected binding detected",
})

events, err := audit.ReadAll("/var/log/portwatch/audit.log")
```
