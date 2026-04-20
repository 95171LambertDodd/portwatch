# Schema Validator

The `schema` package validates observed port entries against operator-defined rules.

## Rules

Each `Rule` specifies:

| Field       | Type     | Description                              |
|-------------|----------|------------------------------------------|
| `MinPort`   | uint16   | Lowest acceptable port number            |
| `MaxPort`   | uint16   | Highest acceptable port number           |
| `Protocols` | []string | Allowed protocols (`tcp`, `udp`); empty means any |

## Example

```go
v, err := schema.New([]schema.Rule{
    {MinPort: 1024, MaxPort: 49151, Protocols: []string{"tcp"}},
    {MinPort: 53,   MaxPort: 53,    Protocols: []string{"udp"}},
})
if err != nil {
    log.Fatal(err)
}

violations := v.ValidateAll(scannedEntries)
for _, e := range violations {
    fmt.Println("schema violation:", e)
}
```

## Integration

The validator is typically wired into the pipeline after scanning and before
alerting, so only schema-compliant entries proceed to baseline checks.

## Error messages

- `port N out of allowed range [min, max]`
- `protocol "X" not in allowed list [...]`
