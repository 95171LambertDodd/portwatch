// Package masking implements field-level redaction of port scan entries.
//
// A Masker is configured with one or more Rules, each targeting a specific
// entry field ("pid", "process", or "addr") and providing a replacement
// string. When applied, the masker returns a copy of the entry with the
// specified fields overwritten, leaving the original untouched.
//
// Example usage:
//
//	m, err := masking.New([]masking.Rule{
//		{Field: "process", Replacement: "[redacted]"},
//		{Field: "pid",     Replacement: "0"},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	masked := m.ApplyAll(entries)
//
// Masking is useful when forwarding alerts to external sinks where
// process names or PIDs should not be disclosed.
package masking
