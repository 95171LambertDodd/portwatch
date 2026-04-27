// Package routing implements rule-based routing of port scan entries to
// named output sinks.
//
// # Overview
//
// A Router holds an ordered list of Rules. Each Rule optionally matches on:
//   - Protocol  ("tcp" or "udp")
//   - Port range (MinPort..MaxPort, inclusive)
//   - Tag        (exact string match)
//
// The first matching rule wins and its Destination label is returned.
//
// A Dispatcher wraps a Router together with a map of named Writer sinks.
// Calling Dispatch routes each PortEntry to the correct sink. Unmatched
// entries are optionally written to a dead-letter io.Writer.
//
// # Example
//
//	rules := []routing.Rule{
//		{Name: "privileged", MaxPort: 1023, Destination: "critical"},
//		{Name: "ephemeral",  MinPort: 32768, Destination: "noise"},
//		{Name: "default",   Destination: "general"},
//	}
//	router, _ := routing.New(rules)
//	dispatcher, _ := routing.NewDispatcher(router, sinks, os.Stderr)
//	dispatcher.Dispatch(entry)
package routing
