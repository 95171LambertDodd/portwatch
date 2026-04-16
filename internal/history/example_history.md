# History Module

The `history` package provides append-only event recording for port binding changes.

## File Format

Events are stored as **JSON Lines** (one JSON object per line) for easy streaming and parsing.

### Example record

```json
{"timestamp":"2024-01-01T00:00:00Z","proto":"tcp","addr":"0.0.0.0","port":8080,"pid":1234,"process":"nginx","kind":"new"}
```

## Event Kinds

| Kind   | Meaning                              |
|--------|--------------------------------------|
| `new`  | A port binding appeared              |
| `gone` | A previously seen binding disappeared|

## Usage

```go
r, err := history.NewRecorder("/var/lib/portwatch/history.jsonl")
if err != nil { log.Fatal(err) }

r.Record(history.Event{
    Timestamp: time.Now(),
    Proto:     "tcp",
    Addr:      "0.0.0.0",
    Port:      9090,
    PID:       42,
    Process:   "myapp",
    Kind:      "new",
})

events, err := history.ReadAll("/var/lib/portwatch/history.jsonl")
```

## Notes

- `NewRecorder` creates intermediate directories automatically.
- `Record` is safe for concurrent use.
- `ReadAll` returns `nil, nil` if the file does not exist.
