# Process Resolver

The `process` package provides a lightweight resolver that maps a PID to its
human-readable name and command-line arguments by reading the Linux `/proc`
filesystem.

## Usage

```go
resolver := process.NewResolver("/proc") // pass "" for default

info, err := resolver.Lookup(pid)
if err != nil {
    log.Printf("could not resolve pid %d: %v", pid, err)
    return
}

fmt.Printf("port bound by: %s (pid %d) — %s\n", info.Name, info.PID, info.Cmdline)
```

## Fields

| Field     | Source             | Description                          |
|-----------|--------------------|--------------------------------------|
| `PID`     | argument           | Numeric process ID                   |
| `Name`    | `/proc/<pid>/comm` | Short process name (≤15 chars)       |
| `Cmdline` | `/proc/<pid>/cmdline` | Full command line, space-separated |

## Notes

- `Cmdline` is best-effort: if the file is missing or unreadable the field is
  left empty and no error is returned.
- `Name` is required; a missing `comm` file returns an error.
- The `procRoot` parameter is injectable for testing with a temporary directory.
