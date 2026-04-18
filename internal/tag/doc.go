// Package tag assigns human-readable labels to port entries based on
// configurable rules. Each rule specifies an optional port, an optional
// protocol constraint ("tcp" or "udp"), and a required label string.
//
// Rules are evaluated in order; all matching rules contribute labels,
// allowing a single entry to carry multiple tags (e.g. "dns", "critical").
//
// Example usage:
//
//	rules := []tag.Rule{
//		{Port: 22,  Protocol: "tcp", Label: "ssh"},
//		{Port: 53,                   Label: "dns"},
//		{Port: 443, Protocol: "tcp", Label: "https"},
//	}
//	tgr, err := tag.New(rules)
//	if err != nil { ... }
//	labels := tgr.Tag(entry) // e.g. ["dns"]
package tag
