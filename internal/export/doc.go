// Package export serialises port scan results to external formats.
//
// Supported formats:
//
//	"json" — pretty-printed JSON array of records.
//	"csv"  — comma-separated values with a header row.
//
// Example usage:
//
//	e, err := export.New("csv")
//	if err != nil { log.Fatal(err) }
//	if err := e.Write(os.Stdout, entries); err != nil { log.Fatal(err) }
package export
