// Package schema provides a rule-based validator for port entries
// observed by the portwatch scanner.
//
// A Validator is constructed with one or more Rule values, each of which
// constrains the acceptable port range and/or protocol of an entry.
//
// Example usage:
//
//	v, err := schema.New([]schema.Rule{
//		{MinPort: 1024, MaxPort: 49151, Protocols: []string{"tcp"}},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := v.Validate(entry); err != nil {
//		fmt.Println("violation:", err)
//	}
//
// ValidateAll may be used to check an entire snapshot at once,
// returning one error per violating entry.
package schema
